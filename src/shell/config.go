package shell

import (
	"flag"
	"fmt"
	"io"
)

// notSet выводится вместо значения параметра, который не был задан.
const notSet = "(не задан)"

// Config — параметры запуска эмулятора, заданные в командной строке.
type Config struct {
	// VFSPath — путь к физическому расположению VFS на диске.
	// Используется начиная с этапа 3.
	VFSPath string
	// Prompt — пользовательское приглашение к вводу. Если пустое,
	// приглашение собирается из реальных данных ОС.
	Prompt string
	// ScriptPath — путь к стартовому скрипту с командами эмулятора.
	ScriptPath string
}

// Pair — один параметр эмулятора в формате ключ-значение.
type Pair struct {
	Key   string
	Value string
}

// ParseFlags разбирает аргументы командной строки и возвращает параметры
// запуска. При неизвестном флаге возвращается ошибка.
func ParseFlags(args []string, out io.Writer) (Config, error) {
	var cfg Config
	set := flag.NewFlagSet("emulator", flag.ContinueOnError)
	set.SetOutput(out)
	set.StringVar(&cfg.VFSPath, "vfs", "", "путь к физическому расположению VFS")
	set.StringVar(&cfg.Prompt, "prompt", "", "приглашение к вводу в REPL")
	set.StringVar(&cfg.ScriptPath, "script", "", "путь к стартовому скрипту")
	if err := set.Parse(args); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Pairs возвращает все параметры эмулятора в формате ключ-значение.
// Используется и отладочным выводом при запуске, и командой conf-dump.
func (c Config) Pairs() []Pair {
	return []Pair{
		{Key: "vfs", Value: valueOrNotSet(c.VFSPath)},
		{Key: "prompt", Value: valueOrNotSet(c.Prompt)},
		{Key: "script", Value: valueOrNotSet(c.ScriptPath)},
	}
}

// PromptString возвращает приглашение к вводу: заданное параметром
// или собранное из реальных данных ОС, если параметр не указан.
func (c Config) PromptString() string {
	if c.Prompt == "" {
		return Prompt()
	}
	return c.Prompt + " "
}

// valueOrNotSet подставляет пометку вместо пустого значения параметра.
func valueOrNotSet(value string) string {
	if value == "" {
		return notSet
	}
	return value
}

// PrintDebug печатает отладочный список всех параметров эмулятора.
// Выводится один раз при запуске, как требует этап «Конфигурация».
func PrintDebug(out io.Writer, cfg Config) {
	fmt.Fprintln(out, "[debug] параметры запуска эмулятора:")
	for _, pair := range cfg.Pairs() {
		fmt.Fprintf(out, "[debug]   %s = %s\n", pair.Key, pair.Value)
	}
}
