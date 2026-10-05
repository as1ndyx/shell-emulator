package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestTouchCreatesEmptyFile(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "touch new.txt")

	if !strings.Contains(mustRun(t, s, "ls"), "new.txt") {
		t.Error("файл не появился в каталоге")
	}
	if out := mustRun(t, s, "cat new.txt"); out != "" {
		t.Errorf("новый файл должен быть пустым, получили %q", out)
	}
}

func TestTouchSeveralFilesAndNestedPath(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "touch docs/notes/idea.txt a.txt")

	if out := mustRun(t, s, "ls docs/notes"); out != "draft.txt  idea.txt  todo.txt" {
		t.Errorf("вложенный каталог: %q", out)
	}
	if !strings.Contains(mustRun(t, s, "ls"), "a.txt") {
		t.Error("второй файл не создан")
	}
}

func TestTouchExistingFileKeepsContent(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "touch readme.md")

	if out := mustRun(t, s, "cat readme.md"); out != "Проект с вложенными каталогами." {
		t.Errorf("содержимое изменилось: %q", out)
	}
}

func TestTouchErrors(t *testing.T) {
	s := loadSession(t, "deep")
	for _, line := range []string{"touch", "touch nosuchdir/f.txt", "touch readme.md/f.txt"} {
		if _, err := run(s, line); err == nil {
			t.Errorf("%q: ожидали ошибку", line)
		}
	}
	if _, err := run(shell.NewSession(shell.Config{}, nil), "touch a.txt"); err == nil {
		t.Error("без VFS ожидали ошибку")
	}
}

func TestCpToNewFile(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "cp readme.md copy.md")

	if out := mustRun(t, s, "cat copy.md"); out != "Проект с вложенными каталогами." {
		t.Errorf("копия: %q", out)
	}
}

func TestCpIntoDirectory(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "cp readme.md docs")

	if out := mustRun(t, s, "ls docs"); out != "notes/  overview.txt  readme.md" {
		t.Errorf("файл не скопирован в каталог: %q", out)
	}
}

func TestCpOverwritesExistingFile(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "cp etc/config.ini docs/notes/todo.txt")

	if out := mustRun(t, s, "cat docs/notes/todo.txt"); out != "Настройки приложения." {
		t.Errorf("файл не перезаписан: %q", out)
	}
}

func TestCpRelativeToCurrentDirectory(t *testing.T) {
	s := loadSession(t, "deep")
	mustRun(t, s, "cd docs")

	mustRun(t, s, "cp ../etc/config.ini settings.ini")

	if !strings.Contains(mustRun(t, s, "ls"), "settings.ini") {
		t.Error("копия не появилась в текущем каталоге")
	}
}

// TestCpCopyIsIndependent проверяет, что после перезаписи копии
// оригинал остаётся прежним.
func TestCpCopyIsIndependent(t *testing.T) {
	s := loadSession(t, "deep")
	mustRun(t, s, "cp readme.md copy.md")

	mustRun(t, s, "cp etc/config.ini copy.md")

	if out := mustRun(t, s, "cat readme.md"); out != "Проект с вложенными каталогами." {
		t.Errorf("оригинал изменился вместе с копией: %q", out)
	}
}

func TestCpErrors(t *testing.T) {
	s := loadSession(t, "deep")
	lines := []string{
		"cp", "cp readme.md", "cp a b c", "cp docs backup", "cp nosuch.txt c.txt",
		"cp readme.md nosuchdir/c.md", "cp readme.md readme.md", "cp readme.md .",
	}
	for _, line := range lines {
		if _, err := run(s, line); err == nil {
			t.Errorf("%q: ожидали ошибку", line)
		}
	}
}

// TestChangesStayInMemory проверяет, что touch и cp меняют VFS только
// в памяти: повторная загрузка с диска даёт исходное дерево.
func TestChangesStayInMemory(t *testing.T) {
	s := loadSession(t, "deep")
	before := s.VFS.Tree()

	mustRun(t, s, "touch new.txt")
	mustRun(t, s, "cp readme.md docs")

	if loadSession(t, "deep").VFS.Tree() != before {
		t.Error("изменения попали на диск")
	}
	if s.VFS.Tree() == before {
		t.Error("изменения не видны в памяти")
	}
}
