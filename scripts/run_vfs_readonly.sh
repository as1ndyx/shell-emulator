#!/bin/sh
# Проверка требования «все операции производятся в памяти»:
# контрольные суммы файлов VFS до и после работы эмулятора совпадают.
set -e
cd "$(dirname "$0")/.."

before=$(find vfs/deep -type f -exec shasum {} \; | sort)

printf 'vfs-tree\nls\nexit\n' | go run ./src -vfs vfs/deep > /dev/null

after=$(find vfs/deep -type f -exec shasum {} \; | sort)

if [ "$before" = "$after" ]; then
	echo "VFS на диске не изменилась: все операции выполнены в памяти"
else
	echo "VFS на диске изменилась"
	exit 1
fi
