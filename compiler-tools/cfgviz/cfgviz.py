#!/usr/bin/env python3
"""
CFG Visualization Tool

This script generates visualizations of Control Flow Graphs (CFGs) from Ballerina source files.
It orchestrates the Ballerina interpreter (to generate DOT format) and Graphviz (to render images).

Usage:
    python3 cfgviz.py <source-file.bal> [options]
    python3 cfgviz.py corpus/cfg/subset1/01-loop/while02-v.bal
    python3 cfgviz.py myfile.bal -f svg -o output.svg
    python3 cfgviz.py myfile.bal --no-open

Requirements:
    - Graphviz (dot command) must be installed
    - Go must be installed (to run the interpreter)
"""

import argparse
import os
import platform
import subprocess
import sys
import tempfile
from pathlib import Path


def check_graphviz():
    """Check if Graphviz (dot command) is available."""
    try:
        subprocess.run(["dot", "-V"], capture_output=True, check=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False


def print_graphviz_install_help():
    """Print installation instructions for Graphviz."""
    system = platform.system()
    print("Error: Graphviz is not installed or not in PATH.", file=sys.stderr)
    print("\nTo install Graphviz:", file=sys.stderr)

    if system == "Darwin":
        print("  brew install graphviz", file=sys.stderr)
    elif system == "Linux":
        print("  sudo apt install graphviz  # Debian/Ubuntu", file=sys.stderr)
        print("  sudo dnf install graphviz  # Fedora", file=sys.stderr)
    elif system == "Windows":
        print("  choco install graphviz  # Chocolatey", file=sys.stderr)
        print("  Or download from: https://graphviz.org/download/", file=sys.stderr)
    else:
        print("  Visit: https://graphviz.org/download/", file=sys.stderr)


def find_project_root():
    """Find the project root directory (where go.mod is located)."""
    # Script is in compiler-tools/cfgviz/ directory, project root is two levels up
    script_dir = Path(__file__).resolve().parent
    project_root = script_dir.parent.parent

    # Verify go.mod exists
    if not (project_root / "go.mod").exists():
        print(f"Error: go.mod not found in {project_root}", file=sys.stderr)
        sys.exit(1)

    return project_root


def run_interpreter(project_root, source_file):
    """
    Run the Ballerina interpreter with --dump-cfg --format=dot.

    Returns:
        tuple: (stdout, stderr, returncode)
    """
    # Convert source file to absolute path
    source_path = Path(source_file).resolve()

    if not source_path.exists():
        print(f"Error: Source file not found: {source_file}", file=sys.stderr)
        sys.exit(1)

    # Run from project root
    cmd = [
        "go",
        "run",
        "./cli/cmd",
        "run",
        "--dump-cfg",
        "--format=dot",
        str(source_path),
    ]

    print(f"Running interpreter: {' '.join(cmd)}", file=sys.stderr)

    try:
        result = subprocess.run(cmd, cwd=project_root, capture_output=True, text=True)
        return result.stdout, result.stderr, result.returncode
    except Exception as e:
        print(f"Error running interpreter: {e}", file=sys.stderr)
        sys.exit(1)


def extract_dot_content(stdout, stderr):
    """
    Extract DOT content from interpreter output.

    The DOT content is in stdout, between BEGIN CFG and END CFG markers in stderr.
    We need to check if the CFG was generated (markers present in stderr).
    The program output may also be in stdout after the DOT content, so we need
    to extract only the DOT graph part.

    Returns:
        str: DOT content, or None if not found
    """
    # Check if CFG was generated (look for markers in stderr)
    if "BEGIN CFG" not in stderr or "END CFG" not in stderr:
        print("Error: No CFG output found in interpreter output", file=sys.stderr)
        return None

    # DOT content is in stdout, but program output may also be there
    # Extract from "digraph CFG {" to the closing "}"
    stdout = stdout.strip()

    if not stdout or not stdout.startswith("digraph"):
        print("Error: Invalid DOT content in stdout", file=sys.stderr)
        print(f"stdout preview: {stdout[:200]}", file=sys.stderr)
        return None

    # Find the DOT graph boundaries
    # Count braces to find the matching closing brace
    brace_count = 0
    in_graph = False
    end_pos = 0

    for i, char in enumerate(stdout):
        if char == "{":
            brace_count += 1
            in_graph = True
        elif char == "}":
            brace_count -= 1
            if in_graph and brace_count == 0:
                end_pos = i + 1
                break

    if end_pos == 0:
        print("Error: Could not find complete DOT graph", file=sys.stderr)
        return None

    dot_content = stdout[:end_pos].strip()
    return dot_content


def render_dot(dot_content, output_file, output_format):
    """
    Render DOT content to an image using Graphviz.

    Args:
        dot_content: DOT format string
        output_file: Output file path
        output_format: Image format (png, svg, pdf, etc.)
    """
    with tempfile.NamedTemporaryFile(mode="w", suffix=".dot", delete=False) as f:
        temp_dot_file = f.name
        f.write(dot_content)

    try:
        cmd = ["dot", f"-T{output_format}", temp_dot_file, "-o", str(output_file)]
        print(f"Generating {output_format.upper()}: {output_file}", file=sys.stderr)

        result = subprocess.run(cmd, capture_output=True, text=True)

        if result.returncode != 0:
            print(f"Error running dot: {result.stderr}", file=sys.stderr)
            sys.exit(1)

        print(f"Success! Generated: {output_file}", file=sys.stderr)
    finally:
        # Clean up temp file
        try:
            os.unlink(temp_dot_file)
        except:
            pass


def open_file(file_path):
    """Open a file with the default application (platform-specific)."""
    system = platform.system()

    try:
        if system == "Darwin":
            subprocess.run(["open", str(file_path)], check=True)
        elif system == "Linux":
            subprocess.run(["xdg-open", str(file_path)], check=True)
        elif system == "Windows":
            os.startfile(str(file_path))
        else:
            print(f"Cannot auto-open file on {system}", file=sys.stderr)
    except Exception as e:
        print(f"Warning: Could not open file: {e}", file=sys.stderr)


def main():
    parser = argparse.ArgumentParser(
        description="Generate CFG visualizations from Ballerina source files",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s myfile.bal
  %(prog)s myfile.bal -f svg
  %(prog)s myfile.bal -o output.png --no-open
  %(prog)s corpus/cfg/subset1/01-loop/while02-v.bal
        """,
    )

    parser.add_argument("source_file", help="Ballerina source file (.bal)")

    parser.add_argument(
        "-f",
        "--format",
        default="png",
        choices=["png", "svg", "pdf", "jpg"],
        help="Output format (default: png)",
    )

    parser.add_argument(
        "-o", "--output", help="Output file path (default: <source>.cfg.<format>)"
    )

    parser.add_argument(
        "--no-open",
        action="store_true",
        help="Don't automatically open the generated file",
    )

    args = parser.parse_args()

    # Check prerequisites
    if not check_graphviz():
        print_graphviz_install_help()
        sys.exit(1)

    # Find project root
    project_root = find_project_root()
    print(f"Project root: {project_root}", file=sys.stderr)

    # Determine output file
    if args.output:
        output_file = Path(args.output)
    else:
        source_path = Path(args.source_file)
        output_file = source_path.parent / f"{source_path.stem}.cfg.{args.format}"

    # Run interpreter
    stdout, stderr, returncode = run_interpreter(project_root, args.source_file)

    # Print stderr (compiler messages)
    if stderr:
        print(stderr, file=sys.stderr)

    # Extract DOT content (proceed even if interpreter had errors)
    dot_content = extract_dot_content(stdout, stderr)

    if not dot_content:
        print("Error: Could not extract DOT content", file=sys.stderr)
        sys.exit(1)

    # Render to image
    render_dot(dot_content, output_file, args.format)

    # Open file if requested
    if not args.no_open:
        open_file(output_file)

    return 0


if __name__ == "__main__":
    sys.exit(main())
