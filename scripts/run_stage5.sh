#!/bin/sh
# Этап 5: все режимы команд touch и cp, включая работу с VFS
# и обработку ошибок. Изменения происходят только в памяти.
set -e
cd "$(dirname "$0")/.."
go run ./src -vfs vfs/deep -script examples/stage5.txt < /dev/null
