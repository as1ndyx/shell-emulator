#!/bin/sh
# Проверка параметра -prompt: приглашение задаётся пользователем.
set -e
cd "$(dirname "$0")/.."
printf 'conf-dump\nuptime\nexit\n' | go run ./src -prompt "myshell$"
