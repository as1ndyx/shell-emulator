#!/bin/sh
# VFS из нескольких файлов без вложенных каталогов.
set -e
cd "$(dirname "$0")/.."
printf 'vfs-tree\nexit\n' | go run ./src -vfs vfs/flat -prompt "flat>"
