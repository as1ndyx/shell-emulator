package shell

import (
	"fmt"
	"strings"
	"time"
)

// Session — состояние работающего эмулятора. До этапа 4 хватало
// параметров запуска, но команда cd меняет текущий каталог, а uptime
// считает время работы. Эти данные должны сохраняться между командами,
// поэтому собраны в одну структуру, которая передаётся по указателю.
type Session struct {
	// Config — параметры запуска из командной строки.
	Config Config
	// VFS — корень загруженной файловой системы или nil, если параметр
	// -vfs не задан.
	VFS *Node
	// Cwd — текущий каталог в виде списка имён от корня VFS:
	// пустой список — корень, ["docs", "notes"] — каталог ~/docs/notes.
	Cwd []string
	// Started — момент запуска эмулятора, от него считает uptime.
	Started time.Time
}

// NewSession создаёт сессию эмулятора. Текущим каталогом становится
// корень VFS, временем запуска — текущий момент.
func NewSession(cfg Config, vfs *Node) *Session {
	return &Session{Config: cfg, VFS: vfs, Started: time.Now()}
}

// PromptString возвращает приглашение к вводу: заданное параметром
// -prompt либо собранное из данных ОС вместе с текущим каталогом.
func (s *Session) PromptString() string {
	if s.Config.Prompt != "" {
		return s.Config.Prompt + " "
	}
	return PromptIn(s.CwdString())
}

// CwdString возвращает текущий каталог в виде ~ или ~/docs/notes.
func (s *Session) CwdString() string {
	if len(s.Cwd) == 0 {
		return homeDir
	}
	return homeDir + "/" + strings.Join(s.Cwd, "/")
}

// resolve превращает путь из команды в путь от корня VFS. Понимает
// абсолютные пути (/etc), пути от корня (~/docs), относительные пути
// (notes/todo.txt), а также «.» и «..».
func (s *Session) resolve(path string) []string {
	parts, rest := s.startOf(path)
	for _, part := range strings.Split(rest, "/") {
		parts = step(parts, part)
	}
	return parts
}

// startOf определяет, откуда отсчитывать путь: от корня VFS или от
// текущего каталога. Возвращает начальный список имён и остаток пути.
// Текущий каталог копируется, а не используется напрямую: иначе разбор
// пути мог бы незаметно испортить состояние сессии.
func (s *Session) startOf(path string) ([]string, string) {
	if path == homeDir || strings.HasPrefix(path, homeDir+"/") {
		return nil, strings.TrimPrefix(path, homeDir)
	}
	if strings.HasPrefix(path, "/") {
		return nil, path
	}
	return append([]string(nil), s.Cwd...), path
}

// step применяет к пути одну его часть: «..» поднимает на уровень выше,
// но не выше корня, как и в настоящей ОС; «.» и пустая часть (из «//»)
// ничего не меняют; остальное — имя вложенного файла или каталога.
func step(parts []string, part string) []string {
	switch part {
	case "", ".":
		return parts
	case "..":
		if len(parts) == 0 {
			return parts
		}
		return parts[:len(parts)-1]
	default:
		return append(parts, part)
	}
}

// requireVFS возвращает ошибку, если VFS не загружена: командам,
// работающим с файлами, без неё делать нечего.
func (s *Session) requireVFS(name string) error {
	if s.VFS == nil {
		return fmt.Errorf("%s: VFS не загружена, укажите параметр -vfs", name)
	}
	return nil
}

// lookup находит файл или каталог по пути из команды. Если его нет,
// возвращает ошибку в стиле UNIX: «имя: путь: нет такого файла».
func (s *Session) lookup(name, path string) (*Node, error) {
	if err := s.requireVFS(name); err != nil {
		return nil, err
	}
	node := s.VFS.find(s.resolve(path))
	if node == nil {
		return nil, fmt.Errorf("%s: %s: нет такого файла или каталога", name, path)
	}
	return node, nil
}
