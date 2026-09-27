package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestRunProcessesSessionAndExits(t *testing.T) {
	in := strings.NewReader("ls -l\nnope\nexit\nls\n")
	out := &strings.Builder{}

	shell.Run(in, out)

	text := out.String()
	if !strings.Contains(text, "ls: -l") {
		t.Errorf("нет вывода команды ls: %q", text)
	}
	if !strings.Contains(text, "ошибка: nope: команда не найдена") {
		t.Errorf("нет сообщения об ошибке: %q", text)
	}
	if strings.Contains(text, "ls: аргументов нет") {
		t.Errorf("команда после exit была выполнена: %q", text)
	}
}

func TestRunPrintsPrompt(t *testing.T) {
	out := &strings.Builder{}

	shell.Run(strings.NewReader("exit\n"), out)

	if !strings.HasPrefix(out.String(), shell.Prompt()) {
		t.Errorf("приглашение не выведено: %q", out.String())
	}
}
