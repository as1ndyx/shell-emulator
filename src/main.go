// Команда emulator запускает REPL эмулятора языка оболочки ОС.
package main

import (
	"fmt"
	"os"

	"shellemu/src/shell"
)

// exitCodeError — код возврата при ошибке запуска эмулятора.
const exitCodeError = 1

// main разбирает параметры командной строки и запускает эмулятор.
// Параметры берутся из os.Args начиная со второго элемента: первый —
// имя самой программы. При ошибке запуска программа завершается
// с кодом возврата 1.
func main() {
	cfg, err := shell.ParseFlags(os.Args[1:], os.Stderr)
	if err != nil {
		os.Exit(exitCodeError)
	}

	if err := shell.Start(cfg, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(exitCodeError)
	}
}
