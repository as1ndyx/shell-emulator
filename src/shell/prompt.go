package shell

import (
	"os"
	"os/user"
)

// Значения, подставляемые вместо реальных данных ОС, если система
// по какой-то причине их не сообщила.
const (
	unknownUser = "unknown"
	unknownHost = "localhost"
)

// homeDir — обозначение корня VFS в приглашении и в путях. Корень
// играет роль домашнего каталога, который в настоящей оболочке
// обозначается тильдой.
const homeDir = "~"

// Prompt возвращает приглашение к вводу вида username@hostname:~$,
// собранное из реальных данных операционной системы.
func Prompt() string {
	return PromptIn(homeDir)
}

// PromptIn возвращает приглашение для указанного текущего каталога,
// например username@hostname:~/docs$, как в настоящей оболочке.
func PromptIn(dir string) string {
	return UserName() + "@" + HostName() + ":" + dir + "$ "
}

// UserName возвращает имя текущего пользователя ОС. Имя берётся
// у системы, а не задаётся в коде: по условию приглашение должно
// строиться из реальных данных ОС. Если имя получить не удалось,
// возвращается "unknown".
func UserName() string {
	current, err := user.Current()
	if err != nil {
		return unknownUser
	}
	return current.Username
}

// HostName возвращает сетевое имя машины, на которой работает эмулятор.
// Если имя получить не удалось, возвращается "localhost".
func HostName() string {
	name, err := os.Hostname()
	if err != nil {
		return unknownHost
	}
	return name
}
