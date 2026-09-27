#!/bin/sh
# Запуск эмулятора языка оболочки ОС.
set -e
cd "$(dirname "$0")"
go run ./src "$@"
