package shell

import (
	"fmt"
	"io"
	"os"
)

// Start запускает эмулятор: печатает отладочный список параметров,
// загружает VFS, создаёт сессию, выполняет стартовый скрипт, если он
// задан, и переходит в интерактивный режим. Сессия одна на всё время
// работы, поэтому скрипт и интерактивный режим видят один и тот же
// текущий каталог. Если в скрипте встретилась команда exit,
// в интерактивный режим эмулятор не переходит. Ошибка возвращается, если
// не удалось загрузить VFS или открыть стартовый скрипт — ошибки самих
// команд обрабатываются внутри цикла и работу не прерывают.
func Start(cfg Config, in io.Reader, out io.Writer) error {
	PrintDebug(out, cfg)

	vfs, err := loadVFSIfSet(cfg, out)
	if err != nil {
		return err
	}
	session := NewSession(cfg, vfs)

	if cfg.ScriptPath != "" {
		done, err := runScript(session, out)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	Run(session, in, out, false)
	return nil
}

// loadVFSIfSet загружает VFS, если задан параметр -vfs. Без параметра
// эмулятор работает без виртуальной файловой системы.
func loadVFSIfSet(cfg Config, out io.Writer) (*Node, error) {
	if cfg.VFSPath == "" {
		return nil, nil
	}
	vfs, err := LoadVFS(cfg.VFSPath)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(out, "[debug] VFS загружена в память из %s\n", cfg.VFSPath)
	return vfs, nil
}

// runScript выполняет команды из стартового скрипта. Файл подставляется
// в REPL вместо клавиатуры, поэтому команды выполняются последовательно,
// а ошибочные строки пропускаются с сообщением. Эхо включено: без него
// на экране был бы виден только вывод, без самого диалога.
// Возвращает true, если скрипт завершился командой exit.
func runScript(s *Session, out io.Writer) (bool, error) {
	file, err := os.Open(s.Config.ScriptPath)
	if err != nil {
		return false, fmt.Errorf("стартовый скрипт не открыт: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(out, "[debug] выполняется скрипт %s\n", s.Config.ScriptPath)
	return Run(s, file, out, true), nil
}
