#!/bin/sh
# Минимальная VFS: один файл в корне.
set -e
cd "$(dirname "$0")/.."
printf 'vfs-tree\nconf-dump\nexit\n' | go run ./src -vfs vfs/minimal
