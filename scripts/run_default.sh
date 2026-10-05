#!/bin/sh
# Запуск без параметров: приглашение собирается из данных ОС,
# стартовый скрипт не выполняется, VFS не загружена.
set -e
cd "$(dirname "$0")/.."
printf 'conf-dump\nuptime\nls\nexit\n' | go run ./src
