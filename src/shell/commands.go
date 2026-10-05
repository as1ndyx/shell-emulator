package shell

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrExit возвращается командой exit и означает запрос на завершение работы.
// Это не настоящая ошибка, а признак выхода: обработчик команды возвращает
// ровно два значения, поэтому сигнал передаётся по тому же каналу,
// что и ошибки, а вызывающий код узнаёт его через errors.Is.
var ErrExit = errors.New("выход из эмулятора")

// handler — функция, выполняющая одну команду эмулятора: получает
// сессию и аргументы, возвращает текст вывода или ошибку.
type handler func(s *Session, args []string) (string, error)

// commands — таблица всех команд эмулятора: имя → обработчик.
// До этапа 4 команды выбирались через switch, но с ростом их числа
// такая функция превысила бы допустимую цикломатическую сложность.
// В таблице добавление команды — это одна новая строка.
var commands = map[string]handler{
	"ls":        cmdLs,
	"cd":        cmdCd,
	"cat":       cmdCat,
	"rev":       cmdRev,
	"uptime":    cmdUptime,
	"conf-dump": cmdConfDump,
	"vfs-tree":  cmdVfsTree,
	"exit":      cmdExit,
}

// Execute выполняет разобранную команду и возвращает текст её вывода.
// Пустой ввод ничего не печатает и ошибкой не считается. Для неизвестной
// команды возвращается ошибка.
func (s *Session) Execute(cmd Command) (string, error) {
	if cmd.Name == "" {
		return "", nil
	}
	run, ok := commands[cmd.Name]
	if !ok {
		return "", fmt.Errorf("%s: команда не найдена", cmd.Name)
	}
	return run(s, cmd.Args)
}

// noArgs возвращает ошибку, если команде, которая не принимает
// аргументов, их всё же передали.
func noArgs(name string, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("%s: аргументы не поддерживаются", name)
	}
	return nil
}

// cmdConfDump — служебная команда: печатает параметры эмулятора
// в формате ключ-значение.
func cmdConfDump(s *Session, args []string) (string, error) {
	if err := noArgs("conf-dump", args); err != nil {
		return "", err
	}
	pairs := s.Config.Pairs()
	lines := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		lines = append(lines, pair.Key+" = "+pair.Value)
	}
	return strings.Join(lines, "\n"), nil
}

// cmdVfsTree — служебная команда: печатает дерево загруженной в память
// VFS. Сама VFS при этом не изменяется.
func cmdVfsTree(s *Session, args []string) (string, error) {
	if err := noArgs("vfs-tree", args); err != nil {
		return "", err
	}
	if err := s.requireVFS("vfs-tree"); err != nil {
		return "", err
	}
	return s.VFS.Tree(), nil
}

// cmdUptime печатает текущее время и сколько работает эмулятор,
// по образцу UNIX-команды uptime: «19:30:12 up 0:02:15».
func cmdUptime(s *Session, args []string) (string, error) {
	if err := noArgs("uptime", args); err != nil {
		return "", err
	}
	now := time.Now()
	return now.Format("15:04:05") + " up " + formatDuration(now.Sub(s.Started)), nil
}

// formatDuration записывает длительность в виде часы:минуты:секунды.
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, d/time.Second)
}

// cmdExit завершает работу эмулятора. Лишние аргументы — ошибка
// пользователя, но не повод закрываться: сообщение будет напечатано,
// а сессия продолжится.
func cmdExit(_ *Session, args []string) (string, error) {
	if err := noArgs("exit", args); err != nil {
		return "", err
	}
	return "", ErrExit
}
