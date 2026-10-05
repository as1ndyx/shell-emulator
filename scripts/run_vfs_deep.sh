#!/bin/sh
# VFS с вложенностью не менее трёх уровней файлов и папок.
set -e
cd "$(dirname "$0")/.."
go run ./src -vfs vfs/deep -prompt "deep>" -script examples/full-test.txt < /dev/null
