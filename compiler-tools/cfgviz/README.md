# CFG Visualization Tool

A tool for generating visual representations of Control Flow Graphs (CFGs) from Ballerina source files.

## Overview

The CFG Visualization tool orchestrates the Ballerina interpreter (to generate CFG data in DOT format) and Graphviz (to render images), making it easy to visualize the control flow of Ballerina programs.

## Prerequisites

Before using this tool, ensure you have the following installed:

1. **Graphviz** - Required for rendering DOT graphs to images
   - macOS: `brew install graphviz`
   - Ubuntu/Debian: `sudo apt install graphviz`
   - Fedora: `sudo dnf install graphviz`
   - Windows: `choco install graphviz` or download from [graphviz.org](https://graphviz.org/download/)

2. **Go** - Required to run the Ballerina interpreter
   - Must be installed and available in your PATH

## Usage

### Basic Usage

Generate a PNG visualization of a Ballerina source file:

```bash
python3 cfgviz.py <source-file.bal>
```

The output file will be created in the same directory as the source file with the name `<filename>.cfg.png`.

### Examples

Generate a PNG (default):
```bash
python3 cfgviz.py myfile.bal
```

Generate an SVG:
```bash
python3 cfgviz.py myfile.bal -f svg
```

Specify a custom output path:
```bash
python3 cfgviz.py myfile.bal -o output.png
```

Generate without auto-opening:
```bash
python3 cfgviz.py myfile.bal --no-open
```

Visualize a test corpus file:
```bash
python3 cfgviz.py corpus/cfg/subset1/01-loop/while02-v.bal
```

## Command-Line Options

- **`source_file`** (required) - Path to the Ballerina source file (.bal)

- **`-f, --format`** - Output image format
  - Choices: `png` (default), `svg`, `pdf`, `jpg`
  - Example: `-f svg`

- **`-o, --output`** - Custom output file path
  - Default: `<source>.cfg.<format>` in the same directory as the source file
  - Example: `-o /path/to/output.png`

- **`--no-open`** - Don't automatically open the generated file
  - By default, the tool opens the generated image with your system's default viewer

## How It Works

1. The script runs the Ballerina interpreter with `--dump-cfg --format=dot` flags
2. The interpreter generates a CFG in DOT format (a graph description language)
3. The DOT content is extracted from the interpreter output
4. Graphviz's `dot` command renders the DOT graph to the specified image format
5. The generated image is optionally opened with the system's default viewer

## Output Format

The generated CFG visualization shows:
- Basic blocks as nodes
- Control flow transitions as edges
- Loop structures and conditional branches
- Entry and exit points of the program

## Troubleshooting

### "Graphviz is not installed or not in PATH"

Make sure Graphviz is installed and the `dot` command is available in your PATH. Try running:
```bash
dot -V
```

If this fails, install Graphviz using the instructions in the Prerequisites section.

### "No CFG output found in interpreter output"

This can occur if:
- The source file has compilation errors
- The interpreter doesn't support CFG generation for the given construct
- The source file path is incorrect

Check the stderr output for compiler error messages.

### "Could not find complete DOT graph"

The DOT output from the interpreter may be malformed. This is usually an internal error. Check the interpreter output for issues.

## Running from Anywhere

The script automatically finds the project root (where `go.mod` is located) and runs the interpreter from there, so you can invoke it from any directory:

```bash
# From project root
python3 compiler-tools/cfgviz/cfgviz.py test.bal

# From compiler-tools/cfgviz
python3 cfgviz.py ../../test.bal

# From anywhere (using absolute path)
python3 /path/to/compiler-tools/cfgviz/cfgviz.py test.bal
```

## Related Files

- **Source code**: `compiler-tools/cfgviz/cfgviz.py`
- **CFG implementation**: See the `semantics/cfg.go` file in the project
- **Test corpus**: `corpus/cfg/` directory contains test files for CFG generation
