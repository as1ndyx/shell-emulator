#!/bin/sh
# Проверка требования «все операции производятся в памяти»: эмулятор
# выполняет команды, в том числе изменяющие VFS (touch, cp), а затем
# контрольные суммы файлов на диске сравниваются с исходными.
set -e
cd "$(dirname "$0")/.."

before=$(find vfs/deep -type f -exec shasum {} \; | sort)

go run ./src -vfs vfs/deep -script examples/stage5.txt < /dev/null > /dev/null

after=$(find vfs/deep -type f -exec shasum {} \; | sort)

if [ "$before" = "$after" ]; then
	echo "VFS на диске не изменилась: все операции выполнены в памяти"
else
	echo "VFS на диске изменилась"
	exit 1
fi
