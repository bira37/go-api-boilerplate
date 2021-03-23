package db

import "fmt"

type Error struct {
	Code    string
	Message string
}

func NewError(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ErrDBInternal(message string) *Error {
	return NewError("ErrDBInternal", message)
}

func ErrDBNotFound(message string) *Error {
	return NewError("ErrDBNotFound", message)
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}
