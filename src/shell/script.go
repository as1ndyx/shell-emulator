package shell

import (
	"fmt"
	"io"
	"os"
)

// Start запускает эмулятор: печатает отладочный список параметров,
// выполняет стартовый скрипт, если он задан, и переходит в интерактивный
// режим. Если в скрипте встретилась команда exit, в интерактивный режим
// эмулятор не переходит. Ошибка возвращается только если стартовый скрипт
// не удалось открыть — ошибки самих команд обрабатываются внутри цикла.
func Start(cfg Config, in io.Reader, out io.Writer) error {
	PrintDebug(out, cfg)

	if cfg.ScriptPath != "" {
		done, err := runScript(cfg, out)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	Run(cfg, in, out, false)
	return nil
}

// runScript выполняет команды из стартового скрипта. Файл подставляется
// в REPL вместо клавиатуры, поэтому команды выполняются последовательно,
// а ошибочные строки пропускаются с сообщением. Эхо включено: без него
// на экране был бы виден только вывод, без самого диалога.
// Возвращает true, если скрипт завершился командой exit.
func runScript(cfg Config, out io.Writer) (bool, error) {
	file, err := os.Open(cfg.ScriptPath)
	if err != nil {
		return false, fmt.Errorf("стартовый скрипт не открыт: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(out, "[debug] выполняется скрипт %s\n", cfg.ScriptPath)
	return Run(cfg, file, out, true), nil
}
