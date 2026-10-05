#!/bin/sh
# Этап 4: все режимы команд ls, cd, cat, rev и uptime, включая работу
# с VFS и обработку ошибок.
set -e
cd "$(dirname "$0")/.."
go run ./src -vfs vfs/deep -script examples/stage4.txt < /dev/null
