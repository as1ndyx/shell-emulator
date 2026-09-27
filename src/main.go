// Команда emulator запускает REPL эмулятора языка оболочки ОС.
package main

import (
	"os"

	"shellemu/src/shell"
)

// main запускает REPL эмулятора на стандартных потоках ввода и вывода.
func main() {
	shell.Run(os.Stdin, os.Stdout)
}
