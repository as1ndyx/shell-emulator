package shell

import "strings"

// Command — результат разбора строки ввода: имя команды и её аргументы.
// Например, для строки "ls -l /home" получится Name = "ls",
// Args = ["-l", "/home"].
type Command struct {
	// Name — первое слово строки, имя команды.
	Name string
	// Args — все остальные слова, аргументы команды.
	Args []string
}

// Parse разбивает строку ввода на команду и аргументы по пробелам:
// первое слово — имя команды, всё остальное — её аргументы. Лишние
// пробелы отбрасываются, поэтому "  cd   /tmp  " превращается в команду
// cd с аргументом /tmp. Для пустой строки возвращается команда с пустым
// именем: без этой проверки обращение к первому слову вызвало бы панику.
func Parse(line string) Command {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return Command{}
	}

	return Command{Name: fields[0], Args: fields[1:]}
}
