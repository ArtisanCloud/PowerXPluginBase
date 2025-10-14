#!/usr/bin/env python3
import hashlib
import os
import sys


def write_hashes(root: str, output: str) -> None:
    entries = []
    for dirpath, _, filenames in os.walk(root):
        for name in filenames:
            full = os.path.join(dirpath, name)
            rel = os.path.relpath(full, root)
            with open(full, "rb") as fh:
                data = fh.read()
            entries.append((rel, hashlib.sha256(data).hexdigest()))
    entries.sort()
    with open(output, "w", encoding="utf-8") as out:
        for rel, sha in entries:
            out.write(f"{sha}  {rel}\n")


def main(argv: list[str]) -> int:
    if len(argv) != 3:
        print("usage: hash_package.py <root> <output>", file=sys.stderr)
        return 1
    root, output = argv[1], argv[2]
    write_hashes(root, output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
