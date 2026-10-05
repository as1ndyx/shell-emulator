package tests

import (
	"strings"
	"testing"

	"shellemu/src/shell"
)

func TestLoadVFSMinimal(t *testing.T) {
	vfs, err := shell.LoadVFS("../vfs/minimal")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !vfs.IsDir || vfs.Name != "minimal" {
		t.Errorf("корень загружен неверно: %+v", vfs)
	}
	if len(vfs.Children) != 1 || vfs.Children[0].Name != "readme.txt" {
		t.Errorf("содержимое загружено неверно: %+v", vfs.Children)
	}
}

// TestLoadVFSKeepsFileContent проверяет, что файлы попадают в память
// вместе с содержимым.
func TestLoadVFSKeepsFileContent(t *testing.T) {
	vfs, err := shell.LoadVFS("../vfs/flat")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !strings.Contains(string(vfs.Children[0].Data), "первый файл") {
		t.Errorf("содержимое файла не загружено: %q", vfs.Children[0].Data)
	}
}

func TestLoadVFSDeepStructure(t *testing.T) {
	vfs, err := shell.LoadVFS("../vfs/deep")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	tree := vfs.Tree()
	for _, want := range []string{"docs/", "notes/", "todo.txt", "src/", "main.txt"} {
		if !strings.Contains(tree, want) {
			t.Errorf("в дереве нет %q:\n%s", want, tree)
		}
	}
}

func TestLoadVFSMissingPath(t *testing.T) {
	_, err := shell.LoadVFS("../vfs/nosuchdir")
	if err == nil {
		t.Error("ожидали ошибку о ненайденном пути")
	}
}

func TestLoadVFSNotADirectory(t *testing.T) {
	_, err := shell.LoadVFS("../vfs/minimal/readme.txt")
	if err == nil || !strings.Contains(err.Error(), "не является директорией") {
		t.Errorf("ожидали ошибку о неверном формате, получили %v", err)
	}
}

func TestVfsTreeCommand(t *testing.T) {
	vfs, err := shell.LoadVFS("../vfs/minimal")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	out, err := shell.NewSession(shell.Config{}, vfs).Execute(shell.Parse("vfs-tree"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !strings.Contains(out, "readme.txt") {
		t.Errorf("дерево не выведено: %q", out)
	}
}

func TestVfsTreeWithoutVFS(t *testing.T) {
	_, err := shell.NewSession(shell.Config{}, nil).Execute(shell.Parse("vfs-tree"))
	if err == nil {
		t.Error("без загруженной VFS ожидали ошибку")
	}
}

func TestVfsTreeWithArgs(t *testing.T) {
	vfs, _ := shell.LoadVFS("../vfs/minimal")

	_, err := shell.NewSession(shell.Config{}, vfs).Execute(shell.Parse("vfs-tree extra"))
	if err == nil {
		t.Error("ожидали ошибку об аргументах")
	}
}
