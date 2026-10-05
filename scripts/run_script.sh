#!/bin/sh
# Проверка параметра -script: команды берутся из стартового скрипта,
# на экране виден диалог — и ввод, и вывод.
set -e
cd "$(dirname "$0")/.."
go run ./src -script examples/start.txt < /dev/null
