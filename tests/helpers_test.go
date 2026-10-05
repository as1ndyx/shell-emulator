package tests

import (
	"testing"

	"shellemu/src/shell"
)

// loadSession создаёт сессию эмулятора с VFS из каталога vfs проекта.
func loadSession(t *testing.T, name string) *shell.Session {
	t.Helper()
	vfs, err := shell.LoadVFS("../vfs/" + name)
	if err != nil {
		t.Fatalf("VFS %s не загружена: %v", name, err)
	}
	return shell.NewSession(shell.Config{}, vfs)
}

// run выполняет строку ввода в сессии и возвращает вывод и ошибку.
func run(s *shell.Session, line string) (string, error) {
	return s.Execute(shell.Parse(line))
}

// mustRun выполняет строку ввода и завершает тест, если команда
// вернула ошибку.
func mustRun(t *testing.T, s *shell.Session, line string) string {
	t.Helper()
	out, err := run(s, line)
	if err != nil {
		t.Fatalf("%q: неожиданная ошибка: %v", line, err)
	}
	return out
}
