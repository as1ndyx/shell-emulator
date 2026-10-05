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

// loadDir рекурсивно читает содержимое директории в память.
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

	node.sortChildren()
	return node, nil
}

// sortChildren упорядочивает вложенные узлы по имени. Порядок обхода
// директории на диске не гарантирован, поэтому без сортировки вывод
// ls и дерева мог бы отличаться от запуска к запуску.
func (n *Node) sortChildren() {
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
}

// addChild добавляет узел в каталог, сохраняя сортировку по имени.
// Изменяется только дерево в памяти, файлы на диске не затрагиваются.
func (n *Node) addChild(child *Node) {
	n.Children = append(n.Children, child)
	n.sortChildren()
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

// child возвращает вложенный узел каталога по имени или nil, если его
// нет. У файла вложенных узлов не бывает.
func (n *Node) child(name string) *Node {
	if !n.IsDir {
		return nil
	}
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

// find спускается от узла по списку имён и возвращает найденный узел
// или nil, если на каком-то шаге нужного имени нет. Пустой список
// означает сам узел.
func (n *Node) find(parts []string) *Node {
	node := n
	for _, name := range parts {
		node = node.child(name)
		if node == nil {
			return nil
		}
	}
	return node
}

// displayName возвращает имя для вывода: каталоги помечаются косой
// чертой в конце, как в ls -F, чтобы их было видно среди файлов.
func (n *Node) displayName() string {
	if n.IsDir {
		return n.Name + "/"
	}
	return n.Name
}

// Tree возвращает дерево узла в виде текста с отступами.
func (n *Node) Tree() string {
	var builder strings.Builder
	n.writeTree(&builder, "")
	return strings.TrimRight(builder.String(), "\n")
}

// writeTree рекурсивно печатает узел и его потомков с отступом.
func (n *Node) writeTree(builder *strings.Builder, indent string) {
	fmt.Fprintf(builder, "%s%s\n", indent, n.displayName())
	for _, child := range n.Children {
		child.writeTree(builder, indent+"  ")
	}
}
