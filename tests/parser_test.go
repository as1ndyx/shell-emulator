package tests

import (
	"slices"
	"testing"

	"shellemu/src/shell"
)

func TestParseCommandWithArgs(t *testing.T) {
	cmd := shell.Parse("ls -l /home")
	if cmd.Name != "ls" {
		t.Errorf("имя команды: получили %q, ожидали \"ls\"", cmd.Name)
	}
	if !slices.Equal(cmd.Args, []string{"-l", "/home"}) {
		t.Errorf("аргументы: получили %v, ожидали [-l /home]", cmd.Args)
	}
}

func TestParseExtraSpaces(t *testing.T) {
	cmd := shell.Parse("   cd    /tmp   ")
	if cmd.Name != "cd" || len(cmd.Args) != 1 || cmd.Args[0] != "/tmp" {
		t.Errorf("получили %+v, ожидали cd с аргументом /tmp", cmd)
	}
}

func TestParseEmptyLine(t *testing.T) {
	cmd := shell.Parse("     ")
	if cmd.Name != "" || len(cmd.Args) != 0 {
		t.Errorf("пустая строка разобрана как %+v", cmd)
	}
}
