package db

import (
	"errors"
	"fmt"
)

var ErrNotEnoughBalance = errors.New("not enough balance")

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
