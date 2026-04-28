#!/usr/bin/env python3
import pathlib
import sys


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: normalize_coverage_paths.py <profile> <module_path> <module_dir>")

    profile_path = pathlib.Path(sys.argv[1])
    module_path = sys.argv[2]
    module_dir = sys.argv[3].removeprefix("./")
    if module_dir == ".":
        module_dir = ""

    if not profile_path.exists():
        return
    if not module_path or module_path == module_dir:
        return

    prefix = module_path + "/"
    target_prefix = module_dir + "/" if module_dir else ""
    lines = profile_path.read_text(encoding="utf-8").splitlines()
    normalized = []
    for line in lines:
        if line.startswith("mode:"):
            normalized.append(line)
            continue
        if line.startswith(prefix):
            normalized.append(target_prefix + line[len(prefix):])
        else:
            normalized.append(line)
    profile_path.write_text("\n".join(normalized) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
