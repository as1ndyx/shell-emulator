package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestRunProcessesSessionAndExits(t *testing.T) {
	in := strings.NewReader("uptime\nnope\nexit\nconf-dump\n")
	out := &strings.Builder{}

	if !shell.Run(shell.NewSession(shell.Config{}, nil), in, out, false) {
		t.Error("команда exit должна завершать цикл")
	}

	text := out.String()
	if !strings.Contains(text, " up ") {
		t.Errorf("нет вывода команды uptime: %q", text)
	}
	if !strings.Contains(text, "ошибка: nope: команда не найдена") {
		t.Errorf("нет сообщения об ошибке: %q", text)
	}
	if strings.Contains(text, "vfs = ") {
		t.Errorf("команда после exit была выполнена: %q", text)
	}
}

func TestRunPrintsPrompt(t *testing.T) {
	out := &strings.Builder{}

	shell.Run(shell.NewSession(shell.Config{}, nil), strings.NewReader("exit\n"), out, false)

	if !strings.HasPrefix(out.String(), shell.Prompt()) {
		t.Errorf("приглашение не выведено: %q", out.String())
	}
}

func TestRunUsesCustomPrompt(t *testing.T) {
	s := shell.NewSession(shell.Config{Prompt: "myshell$"}, nil)
	out := &strings.Builder{}

	shell.Run(s, strings.NewReader("exit\n"), out, false)

	if !strings.HasPrefix(out.String(), "myshell$ ") {
		t.Errorf("приглашение из параметра не использовано: %q", out.String())
	}
}

// TestRunKeepsDirectoryBetweenCommands проверяет, что второй команде
// виден каталог, выбранный первой.
func TestRunKeepsDirectoryBetweenCommands(t *testing.T) {
	out := &strings.Builder{}

	shell.Run(loadSession(t, "deep"), strings.NewReader("cd docs\nls\n"), out, false)

	text := out.String()
	if !strings.Contains(text, ":~/docs$ ") || !strings.Contains(text, "notes/  overview.txt") {
		t.Errorf("каталог не сохранился между командами: %q", text)
	}
}

// TestRunEchoesInputWhenEnabled проверяет, что при включённом эхе
// на экране видны и команда, и её вывод.
func TestRunEchoesInputWhenEnabled(t *testing.T) {
	out := &strings.Builder{}

	shell.Run(shell.NewSession(shell.Config{}, nil), strings.NewReader("uptime\n"), out, true)

	text := out.String()
	if !strings.Contains(text, "$ uptime\n") {
		t.Errorf("введённая команда не продублирована: %q", text)
	}
	if !strings.Contains(text, " up ") {
		t.Errorf("нет вывода команды: %q", text)
	}
}

// TestRunSkipsBadLinesAndContinues проверяет, что ошибочные строки
// пропускаются: последняя команда после них всё равно выполняется.
func TestRunSkipsBadLinesAndContinues(t *testing.T) {
	in := strings.NewReader("qwerty\nexit extra\nuptime\n")
	out := &strings.Builder{}

	if shell.Run(shell.NewSession(shell.Config{}, nil), in, out, true) {
		t.Error("ошибочные строки не должны завершать работу")
	}

	if !strings.Contains(out.String(), " up ") {
		t.Errorf("выполнение прервалось на ошибке: %q", out.String())
	}
}
