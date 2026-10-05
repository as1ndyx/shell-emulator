#!/bin/sh
# Запуск без параметров: приглашение собирается из данных ОС,
# стартовый скрипт не выполняется.
set -e
cd "$(dirname "$0")/.."
printf 'conf-dump\nls -l\nexit\n' | go run ./src
