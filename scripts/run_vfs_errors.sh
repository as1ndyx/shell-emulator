#!/bin/sh
# Ошибки загрузки VFS: несуществующий путь и неверный формат источника.
cd "$(dirname "$0")/.."

echo "--- путь не найден ---"
go run ./src -vfs vfs/nosuchdir < /dev/null || echo "код возврата: $?"

echo "--- источником указан файл, а не директория ---"
go run ./src -vfs vfs/minimal/readme.txt < /dev/null || echo "код возврата: $?"

echo "--- команда vfs-tree без загруженной VFS ---"
printf 'vfs-tree\nexit\n' | go run ./src
