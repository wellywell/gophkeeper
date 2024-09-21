package db

import (
	"fmt"
)

type UserExistsError struct {
	Username string
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("User %s exists", e.Username)
}

type UserNotFoundError struct {
	Username string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User %s not found", e.Username)
}

// KeyExistsError ошибка повторной записи ключа
type KeyExistsError struct {
	Key string
}

// Error стандартный метод интерфейса error
func (e *KeyExistsError) Error() string {
	return fmt.Sprintf("Key %s already exists", e.Key)
}

// KeyNotFoundError ошибка при попытке достать несуществующих ключ
type KeyNotFoundError struct {
	Key string
}

// Error стандартный метод интерфейса error
func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key %s not found", e.Key)
}
