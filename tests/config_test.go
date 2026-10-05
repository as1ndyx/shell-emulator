package tests

import (
	"io"
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestParseFlagsAllParameters(t *testing.T) {
	args := []string{"-vfs", "/tmp/vfs", "-prompt", "myshell$", "-script", "s.txt"}

	cfg, err := shell.ParseFlags(args, io.Discard)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if cfg.VFSPath != "/tmp/vfs" {
		t.Errorf("vfs: получили %q", cfg.VFSPath)
	}
	if cfg.Prompt != "myshell$" {
		t.Errorf("prompt: получили %q", cfg.Prompt)
	}
	if cfg.ScriptPath != "s.txt" {
		t.Errorf("script: получили %q", cfg.ScriptPath)
	}
}

func TestParseFlagsWithoutParameters(t *testing.T) {
	cfg, err := shell.ParseFlags(nil, io.Discard)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if cfg != (shell.Config{}) {
		t.Errorf("ожидали пустую конфигурацию, получили %+v", cfg)
	}
}

func TestParseFlagsUnknownParameter(t *testing.T) {
	_, err := shell.ParseFlags([]string{"-nosuchflag"}, io.Discard)
	if err == nil {
		t.Error("ожидали ошибку о неизвестном параметре")
	}
}

func TestPromptStringFromConfig(t *testing.T) {
	cfg := shell.Config{Prompt: "myshell$"}

	if cfg.PromptString() != "myshell$ " {
		t.Errorf("получили %q, ожидали \"myshell$ \"", cfg.PromptString())
	}
}

func TestPromptStringFallsBackToOS(t *testing.T) {
	if (shell.Config{}).PromptString() != shell.Prompt() {
		t.Error("без параметра приглашение должно браться из данных ОС")
	}
}

func TestPrintDebugListsAllParameters(t *testing.T) {
	cfg := shell.Config{VFSPath: "/tmp/vfs", ScriptPath: "s.txt"}
	out := &strings.Builder{}

	shell.PrintDebug(out, cfg)

	text := out.String()
	for _, want := range []string{"vfs = /tmp/vfs", "prompt = (не задан)", "script = s.txt"} {
		if !strings.Contains(text, want) {
			t.Errorf("в отладочном выводе нет %q: %s", want, text)
		}
	}
}
