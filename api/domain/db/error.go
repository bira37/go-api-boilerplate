package db

import "fmt"

// Database Errors
const (
	ErrDBInternalError = "ErrDBInternalError"
	ErrDBNotFound      = "ErrDBNotFound"
)

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

func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}
