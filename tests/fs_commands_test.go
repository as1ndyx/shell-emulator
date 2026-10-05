package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestLsCurrentDirectory(t *testing.T) {
	out := mustRun(t, loadSession(t, "deep"), "ls")

	if out != "docs/  etc/  readme.md  src/" {
		t.Errorf("получили %q", out)
	}
}

func TestLsGivenDirectory(t *testing.T) {
	out := mustRun(t, loadSession(t, "deep"), "ls docs/notes")

	if out != "draft.txt  todo.txt" {
		t.Errorf("получили %q", out)
	}
}

func TestLsFilePrintsItsName(t *testing.T) {
	if out := mustRun(t, loadSession(t, "deep"), "ls readme.md"); out != "readme.md" {
		t.Errorf("получили %q", out)
	}
}

func TestLsErrors(t *testing.T) {
	s := loadSession(t, "deep")
	for _, line := range []string{"ls nosuchdir", "ls docs src"} {
		if _, err := run(s, line); err == nil {
			t.Errorf("%q: ожидали ошибку", line)
		}
	}
}

func TestFsCommandsWithoutVFS(t *testing.T) {
	s := shell.NewSession(shell.Config{}, nil)
	for _, line := range []string{"ls", "cd docs", "cat a.txt", "rev a.txt"} {
		_, err := run(s, line)
		if err == nil || !strings.Contains(err.Error(), "VFS не загружена") {
			t.Errorf("%q: ожидали ошибку о VFS, получили %v", line, err)
		}
	}
}

func TestCdChangesDirectoryAndPrompt(t *testing.T) {
	s := loadSession(t, "deep")

	mustRun(t, s, "cd docs/notes")

	if s.CwdString() != "~/docs/notes" {
		t.Errorf("текущий каталог: %q", s.CwdString())
	}
	if !strings.HasSuffix(s.PromptString(), ":~/docs/notes$ ") {
		t.Errorf("каталог не виден в приглашении: %q", s.PromptString())
	}
}

func TestCdRelativeAbsoluteAndParent(t *testing.T) {
	s := loadSession(t, "deep")
	steps := []struct{ line, want string }{
		{"cd docs", "~/docs"},
		{"cd notes", "~/docs/notes"},
		{"cd ..", "~/docs"},
		{"cd ../src/app", "~/src/app"},
		{"cd /etc", "~/etc"},
		{"cd ~/docs", "~/docs"},
		{"cd", "~"},
		{"cd ..", "~"},
	}
	for _, step := range steps {
		mustRun(t, s, step.line)
		if s.CwdString() != step.want {
			t.Errorf("после %q каталог %q, ожидали %q", step.line, s.CwdString(), step.want)
		}
	}
}

// TestCdErrorsKeepDirectory проверяет, что ошибочный cd сообщает
// об ошибке и не меняет текущий каталог.
func TestCdErrorsKeepDirectory(t *testing.T) {
	s := loadSession(t, "deep")
	mustRun(t, s, "cd docs")

	for _, line := range []string{"cd nosuchdir", "cd overview.txt", "cd notes ../etc"} {
		if _, err := run(s, line); err == nil {
			t.Errorf("%q: ожидали ошибку", line)
		}
	}
	if s.CwdString() != "~/docs" {
		t.Errorf("каталог изменился после ошибки: %q", s.CwdString())
	}
}

func TestCatFiles(t *testing.T) {
	s := loadSession(t, "deep")

	if out := mustRun(t, s, "cat docs/notes/todo.txt"); out != "Список задач." {
		t.Errorf("один файл: %q", out)
	}
	if out := mustRun(t, s, "cat docs/notes/todo.txt docs/notes/draft.txt"); out != "Список задач.\nЧерновик." {
		t.Errorf("два файла: %q", out)
	}
}

// TestRevReversesLines проверяет, что rev переворачивает строку,
// не ломая русские буквы.
func TestRevReversesLines(t *testing.T) {
	out := mustRun(t, loadSession(t, "deep"), "rev docs/notes/todo.txt")

	if out != ".чадаз косипС" {
		t.Errorf("получили %q", out)
	}
}

func TestCatAndRevErrors(t *testing.T) {
	s := loadSession(t, "deep")
	for _, line := range []string{"cat", "cat docs", "cat nosuch.txt", "rev", "rev docs", "rev nosuch.txt"} {
		if _, err := run(s, line); err == nil {
			t.Errorf("%q: ожидали ошибку", line)
		}
	}
}
