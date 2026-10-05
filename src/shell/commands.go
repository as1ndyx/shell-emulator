package shell

import (
	"errors"
	"fmt"
	"strings"
)

// ErrExit возвращается командой exit и означает запрос на завершение работы.
// Это не настоящая ошибка, а признак выхода: функция Execute возвращает
// ровно два значения, поэтому сигнал передаётся по тому же каналу,
// что и ошибки, а вызывающий код узнаёт его через errors.Is.
var ErrExit = errors.New("выход из эмулятора")

// Execute выполняет разобранную команду и возвращает текст её вывода.
// Параметры запуска нужны служебной команде conf-dump. Пустой ввод ничего
// не печатает и ошибкой не считается. Команды ls и cd пока заглушки,
// настоящая логика появится на этапе 4. Для неизвестной команды
// возвращается ошибка.
func Execute(cfg Config, cmd Command) (string, error) {
	switch cmd.Name {
	case "":
		return "", nil
	case "ls", "cd":
		return stub(cmd), nil
	case "conf-dump":
		return confDump(cfg, cmd)
	case "exit":
		return exit(cmd)
	default:
		return "", fmt.Errorf("%s: команда не найдена", cmd.Name)
	}
}

// stub — заглушка команды: выводит её имя и полученные аргументы,
// склеенные обратно в строку через пробел. Позволяет убедиться,
// что парсер правильно разобрал строку ввода.
func stub(cmd Command) string {
	if len(cmd.Args) == 0 {
		return cmd.Name + ": аргументов нет"
	}
	return cmd.Name + ": " + strings.Join(cmd.Args, " ")
}

// confDump — служебная команда: печатает параметры эмулятора
// в формате ключ-значение. Аргументы не поддерживаются.
func confDump(cfg Config, cmd Command) (string, error) {
	if len(cmd.Args) > 0 {
		return "", errors.New("conf-dump: аргументы не поддерживаются")
	}
	lines := make([]string, 0, len(cfg.Pairs()))
	for _, pair := range cfg.Pairs() {
		lines = append(lines, pair.Key+" = "+pair.Value)
	}
	return strings.Join(lines, "\n"), nil
}

// exit завершает работу эмулятора. Аргументы не поддерживаются: лишние
// аргументы — ошибка пользователя, но не повод закрываться, поэтому
// сообщение будет напечатано, а сессия продолжится.
func exit(cmd Command) (string, error) {
	if len(cmd.Args) > 0 {
		return "", errors.New("exit: аргументы не поддерживаются")
	}
	return "", ErrExit
}
