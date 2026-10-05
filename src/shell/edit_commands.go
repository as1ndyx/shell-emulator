package shell

import (
	"errors"
	"fmt"
)

// cpOperands — сколько аргументов принимает cp: источник и назначение.
const cpOperands = 2

// cmdTouch создаёт пустые файлы по указанным путям. Существующий файл
// или каталог не изменяется. Все изменения происходят только в памяти:
// исходная директория на диске остаётся прежней.
func cmdTouch(s *Session, args []string) (string, error) {
	if err := s.requireVFS("touch"); err != nil {
		return "", err
	}
	if len(args) == 0 {
		return "", errors.New("touch: укажите файл")
	}
	for _, path := range args {
		if err := s.touch(path); err != nil {
			return "", err
		}
	}
	return "", nil
}

// touch создаёт один пустой файл, если по этому пути ещё ничего нет.
// Настоящий touch у существующего файла обновил бы время изменения,
// но время в VFS не хранится, поэтому существующий файл не меняется.
func (s *Session) touch(path string) error {
	parts := s.resolve(path)
	if s.VFS.find(parts) != nil {
		return nil
	}
	parent, name, err := s.parentOf("touch", path, parts)
	if err != nil {
		return err
	}
	parent.addChild(&Node{Name: name})
	return nil
}

// parentOf находит каталог, в котором должен появиться новый файл,
// и имя этого файла. Ошибка — если такого каталога нет, как в UNIX:
// touch nosuchdir/a.txt не создаёт промежуточные каталоги.
// Список parts не бывает пустым: пустой путь — это корень VFS,
// а он всегда существует, поэтому до вызова дело не доходит.
func (s *Session) parentOf(name, path string, parts []string) (*Node, string, error) {
	last := len(parts) - 1
	parent := s.VFS.find(parts[:last])
	if parent == nil || !parent.IsDir {
		return nil, "", fmt.Errorf("%s: %s: нет такого файла или каталога", name, path)
	}
	return parent, parts[last], nil
}

// cmdCp копирует файл. Если назначение — существующий каталог, копия
// кладётся в него под тем же именем; если существующий файл — его
// содержимое заменяется. Каталоги не копируются, как и в настоящем cp
// без флага -r.
func cmdCp(s *Session, args []string) (string, error) {
	if len(args) != cpOperands {
		return "", errors.New("cp: укажите источник и назначение")
	}
	src, err := s.lookup("cp", args[0])
	if err != nil {
		return "", err
	}
	if src.IsDir {
		return "", fmt.Errorf("cp: %s: это каталог, копирование каталогов не поддерживается", args[0])
	}
	return "", s.copyFile(src, args[0], args[1])
}

// copyFile кладёт копию файла src по пути dst. Если dst — каталог,
// копия ложится внутрь него под тем же именем.
func (s *Session) copyFile(src *Node, srcPath, dst string) error {
	parts := s.resolve(dst)
	target := s.VFS.find(parts)
	if target != nil && target.IsDir {
		parts = append(parts, src.Name)
		target = target.child(src.Name)
	}
	if target == src {
		return fmt.Errorf("cp: %s и %s — один и тот же файл", srcPath, dst)
	}
	if target != nil {
		return overwrite(target, src, dst)
	}
	parent, name, err := s.parentOf("cp", dst, parts)
	if err != nil {
		return err
	}
	parent.addChild(&Node{Name: name, Data: copyBytes(src.Data)})
	return nil
}

// overwrite заменяет содержимое существующего файла target содержимым
// src. Заменить каталог файлом нельзя.
func overwrite(target, src *Node, dst string) error {
	if target.IsDir {
		return fmt.Errorf("cp: %s: нельзя заменить каталог файлом", dst)
	}
	target.Data = copyBytes(src.Data)
	return nil
}

// copyBytes возвращает независимую копию содержимого файла. Без неё
// оригинал и копия делили бы одну область памяти, и изменение одного
// файла меняло бы другой.
func copyBytes(data []byte) []byte {
	return append([]byte(nil), data...)
}
