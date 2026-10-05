#!/bin/sh
# Проверка параметра -vfs: путь к физическому расположению VFS.
set -e
cd "$(dirname "$0")/.."
printf 'conf-dump\nexit\n' | go run ./src -vfs /tmp/vfs
