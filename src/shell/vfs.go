package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Node — узел виртуальной файловой системы: файл или каталог.
// Всё дерево целиком хранится в оперативной памяти, исходная
// директория на диске после загрузки не используется и не изменяется.
type Node struct {
	// Name — имя файла или каталога без пути.
	Name string
	// IsDir — true, если узел является каталогом.
	IsDir bool
	// Data — содержимое файла. У каталога всегда пустое.
	Data []byte
	// Children — вложенные узлы каталога, отсортированные по имени.
	Children []*Node
}

// LoadVFS загружает директорию с диска в память и возвращает корень
// дерева. Источником VFS служит обычная директория пользователя.
// Если путь не существует, возвращается ошибка «файл не найден»; если он
// указывает на одиночный файл, а не на директорию, — ошибка неверного
// формата.
func LoadVFS(root string) (*Node, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("VFS не загружена: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("VFS не загружена: %s не является директорией", root)
	}
	return loadDir(root, filepath.Base(root))
}

// loadDir рекурсивно читает содержимое директории в память. Порядок
// обхода директории на диске не гарантирован, поэтому вложенные узлы
// сортируются по имени: вывод дерева должен быть одинаковым при каждом
// запуске.
func loadDir(path, name string) (*Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("VFS не загружена: %w", err)
	}

	node := &Node{Name: name, IsDir: true}
	for _, entry := range entries {
		child, err := loadEntry(filepath.Join(path, entry.Name()), entry.Name(), entry.IsDir())
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, child)
	}

	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Name < node.Children[j].Name
	})
	return node, nil
}

// loadEntry загружает один элемент директории: каталог рекурсивно,
// файл — вместе с содержимым.
func loadEntry(path, name string, isDir bool) (*Node, error) {
	if isDir {
		return loadDir(path, name)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("VFS не загружена: %w", err)
	}
	return &Node{Name: name, Data: data}, nil
}

// Tree возвращает дерево узла в виде текста с отступами.
func (n *Node) Tree() string {
	var builder strings.Builder
	n.writeTree(&builder, "")
	return strings.TrimRight(builder.String(), "\n")
}

// writeTree рекурсивно печатает узел и его потомков с отступом.
func (n *Node) writeTree(builder *strings.Builder, indent string) {
	suffix := ""
	if n.IsDir {
		suffix = "/"
	}
	fmt.Fprintf(builder, "%s%s%s\n", indent, n.Name, suffix)
	for _, child := range n.Children {
		child.writeTree(builder, indent+"  ")
	}
}
