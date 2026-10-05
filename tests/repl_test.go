package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestRunProcessesSessionAndExits(t *testing.T) {
	in := strings.NewReader("ls -l\nnope\nexit\nls\n")
	out := &strings.Builder{}

	if !shell.Run(shell.Config{}, in, out, false) {
		t.Error("команда exit должна завершать цикл")
	}

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

	shell.Run(shell.Config{}, strings.NewReader("exit\n"), out, false)

	if !strings.HasPrefix(out.String(), shell.Prompt()) {
		t.Errorf("приглашение не выведено: %q", out.String())
	}
}

func TestRunUsesCustomPrompt(t *testing.T) {
	cfg := shell.Config{Prompt: "myshell$"}
	out := &strings.Builder{}

	shell.Run(cfg, strings.NewReader("exit\n"), out, false)

	if !strings.HasPrefix(out.String(), "myshell$ ") {
		t.Errorf("приглашение из параметра не использовано: %q", out.String())
	}
}

// TestRunEchoesInputWhenEnabled проверяет, что при включённом эхе
// на экране видны и команда, и её вывод.
func TestRunEchoesInputWhenEnabled(t *testing.T) {
	out := &strings.Builder{}

	shell.Run(shell.Config{}, strings.NewReader("ls -l\n"), out, true)

	text := out.String()
	if !strings.Contains(text, "$ ls -l\n") {
		t.Errorf("введённая команда не продублирована: %q", text)
	}
	if !strings.Contains(text, "ls: -l") {
		t.Errorf("нет вывода команды: %q", text)
	}
}

// TestRunSkipsBadLinesAndContinues проверяет, что ошибочные строки
// пропускаются: последняя команда после них всё равно выполняется.
func TestRunSkipsBadLinesAndContinues(t *testing.T) {
	in := strings.NewReader("qwerty\nexit extra\nls\n")
	out := &strings.Builder{}

	if shell.Run(shell.Config{}, in, out, true) {
		t.Error("ошибочные строки не должны завершать работу")
	}

	if !strings.Contains(out.String(), "ls: аргументов нет") {
		t.Errorf("выполнение прервалось на ошибке: %q", out.String())
	}
}
