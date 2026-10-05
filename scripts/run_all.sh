#!/bin/sh
# Проверка всех параметров сразу, а также пропуска ошибочных строк
# скрипта и завершения работы по команде exit внутри скрипта.
set -e
cd "$(dirname "$0")/.."

echo "--- все параметры вместе ---"
go run ./src -vfs vfs/deep -prompt "myshell$" -script examples/start.txt < /dev/null

echo "--- ошибочные строки пропускаются ---"
go run ./src -prompt "test>" -script examples/with-errors.txt < /dev/null

echo "--- команда exit внутри скрипта ---"
go run ./src -vfs vfs/deep -script examples/with-exit.txt < /dev/null

echo "--- несуществующий скрипт ---"
go run ./src -script examples/nosuchfile.txt < /dev/null || echo "код возврата: $?"
