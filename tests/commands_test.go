package tests

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestExecuteEmptyLine(t *testing.T) {
	out, err := run(shell.NewSession(shell.Config{}, nil), "   ")
	if out != "" || err != nil {
		t.Errorf("пустой ввод: получили %q, %v", out, err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	_, err := run(shell.NewSession(shell.Config{}, nil), "qwerty")
	if err == nil || !strings.Contains(err.Error(), "команда не найдена") {
		t.Errorf("ожидали ошибку о неизвестной команде, получили %v", err)
	}
}

func TestExecuteExit(t *testing.T) {
	_, err := run(shell.NewSession(shell.Config{}, nil), "exit")
	if !errors.Is(err, shell.ErrExit) {
		t.Errorf("ожидали ErrExit, получили %v", err)
	}
}

func TestExecuteExitWithArgs(t *testing.T) {
	_, err := run(shell.NewSession(shell.Config{}, nil), "exit now")
	if err == nil || errors.Is(err, shell.ErrExit) {
		t.Errorf("ожидали ошибку об аргументах, получили %v", err)
	}
}

func TestExecuteConfDump(t *testing.T) {
	cfg := shell.Config{VFSPath: "/tmp/vfs", Prompt: "sh$", ScriptPath: "s.txt"}

	out := mustRun(t, shell.NewSession(cfg, nil), "conf-dump")

	want := "vfs = /tmp/vfs\nprompt = sh$\nscript = s.txt"
	if out != want {
		t.Errorf("получили %q, ожидали %q", out, want)
	}
}

func TestExecuteConfDumpWithArgs(t *testing.T) {
	if _, err := run(shell.NewSession(shell.Config{}, nil), "conf-dump extra"); err == nil {
		t.Error("ожидали ошибку об аргументах")
	}
}

// TestUptimeFormat проверяет формат вывода uptime: «19:30:12 up 0:00:00».
func TestUptimeFormat(t *testing.T) {
	out := mustRun(t, shell.NewSession(shell.Config{}, nil), "uptime")

	format := regexp.MustCompile(`^\d{2}:\d{2}:\d{2} up \d+:\d{2}:\d{2}$`)
	if !format.MatchString(out) {
		t.Errorf("неверный формат uptime: %q", out)
	}
}

func TestUptimeWithArgs(t *testing.T) {
	if _, err := run(shell.NewSession(shell.Config{}, nil), "uptime extra"); err == nil {
		t.Error("ожидали ошибку об аргументах")
	}
}
