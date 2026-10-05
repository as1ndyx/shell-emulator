package tests

import (
	"errors"
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestExecuteStubPrintsArgs(t *testing.T) {
	out, err := shell.Execute(shell.Config{}, nil, shell.Parse("ls -a /etc"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if out != "ls: -a /etc" {
		t.Errorf("получили %q, ожидали \"ls: -a /etc\"", out)
	}
}

func TestExecuteStubWithoutArgs(t *testing.T) {
	out, err := shell.Execute(shell.Config{}, nil, shell.Parse("cd"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if out != "cd: аргументов нет" {
		t.Errorf("получили %q", out)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	_, err := shell.Execute(shell.Config{}, nil, shell.Parse("qwerty"))
	if err == nil || !strings.Contains(err.Error(), "команда не найдена") {
		t.Errorf("ожидали ошибку о неизвестной команде, получили %v", err)
	}
}

func TestExecuteExit(t *testing.T) {
	_, err := shell.Execute(shell.Config{}, nil, shell.Parse("exit"))
	if !errors.Is(err, shell.ErrExit) {
		t.Errorf("ожидали ErrExit, получили %v", err)
	}
}

func TestExecuteExitWithArgs(t *testing.T) {
	_, err := shell.Execute(shell.Config{}, nil, shell.Parse("exit now"))
	if err == nil || errors.Is(err, shell.ErrExit) {
		t.Errorf("ожидали ошибку об аргументах, получили %v", err)
	}
}

func TestExecuteConfDump(t *testing.T) {
	cfg := shell.Config{VFSPath: "/tmp/vfs", Prompt: "sh$", ScriptPath: "s.txt"}

	out, err := shell.Execute(cfg, nil, shell.Parse("conf-dump"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	want := "vfs = /tmp/vfs\nprompt = sh$\nscript = s.txt"
	if out != want {
		t.Errorf("получили %q, ожидали %q", out, want)
	}
}

func TestExecuteConfDumpWithArgs(t *testing.T) {
	_, err := shell.Execute(shell.Config{}, nil, shell.Parse("conf-dump extra"))
	if err == nil {
		t.Error("ожидали ошибку об аргументах")
	}
}
