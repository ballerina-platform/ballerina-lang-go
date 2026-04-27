#!/usr/bin/env python3
from pathlib import Path


def main() -> None:
    for gomod in sorted(Path(".").rglob("go.mod")):
        path = gomod.as_posix()
        if path == "go.mod":
            continue
        if "/vendor/" in f"/{path}":
            continue
        print(gomod.parent.as_posix())


if __name__ == "__main__":
    main()
