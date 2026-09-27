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

// Prompt возвращает приглашение к вводу вида username@hostname:~$,
// собранное из реальных данных операционной системы.
func Prompt() string {
	return UserName() + "@" + HostName() + ":~$ "
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
