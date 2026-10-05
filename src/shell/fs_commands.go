package shell

import (
	"fmt"
	"strings"
)

// optionalPath достаёт путь из аргументов команды, которая принимает
// не больше одного пути. Если путь не указан, возвращает fallback.
func optionalPath(name string, args []string, fallback string) (string, error) {
	switch len(args) {
	case 0:
		return fallback, nil
	case 1:
		return args[0], nil
	default:
		return "", fmt.Errorf("%s: слишком много аргументов", name)
	}
}

// cmdLs выводит содержимое каталога: без аргументов — текущего,
// с аргументом — указанного. Для файла, как и настоящий ls,
// выводит его имя.
func cmdLs(s *Session, args []string) (string, error) {
	path, err := optionalPath("ls", args, ".")
	if err != nil {
		return "", err
	}
	node, err := s.lookup("ls", path)
	if err != nil {
		return "", err
	}
	if !node.IsDir {
		return path, nil
	}
	names := make([]string, 0, len(node.Children))
	for _, child := range node.Children {
		names = append(names, child.displayName())
	}
	return strings.Join(names, "  "), nil
}

// cmdCd меняет текущий каталог. Без аргументов возвращает в корень VFS,
// так же как cd без аргументов возвращает в домашний каталог. Меняется
// только состояние сессии в памяти, и новый каталог сразу виден
// в приглашении.
func cmdCd(s *Session, args []string) (string, error) {
	path, err := optionalPath("cd", args, homeDir)
	if err != nil {
		return "", err
	}
	node, err := s.lookup("cd", path)
	if err != nil {
		return "", err
	}
	if !node.IsDir {
		return "", fmt.Errorf("cd: %s: не является каталогом", path)
	}
	s.Cwd = s.resolve(path)
	return "", nil
}

// cmdCat выводит содержимое одного или нескольких файлов подряд.
// Перевод строки в конце убирается: его и так добавит REPL при печати.
func cmdCat(s *Session, args []string) (string, error) {
	text, err := readFiles(s, "cat", args)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(text, "\n"), nil
}

// cmdRev выводит строки файлов, перевернув каждую задом наперёд,
// как UNIX-команда rev.
func cmdRev(s *Session, args []string) (string, error) {
	text, err := readFiles(s, "rev", args)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	for i, line := range lines {
		lines[i] = reverse(line)
	}
	return strings.Join(lines, "\n"), nil
}

// readFiles читает содержимое файлов VFS и склеивает его в одну строку.
// Каталог или несуществующий путь — ошибка, как в настоящих cat и rev.
func readFiles(s *Session, name string, paths []string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("%s: укажите файл", name)
	}
	var builder strings.Builder
	for _, path := range paths {
		node, err := s.lookup(name, path)
		if err != nil {
			return "", err
		}
		if node.IsDir {
			return "", fmt.Errorf("%s: %s: это каталог", name, path)
		}
		builder.Write(node.Data)
	}
	return builder.String(), nil
}

// reverse переворачивает строку посимвольно. Работает с рунами,
// а не с байтами: русская буква занимает два байта, и при побайтовом
// перевороте превратилась бы в мусор.
func reverse(line string) string {
	runes := []rune(line)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
