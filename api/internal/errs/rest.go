package errs

import "fmt"

type RestError struct {
	Code       string
	Message    string
	StatusCode int
}

func NewRestError(message string, code string, status int) *RestError {
	return &RestError{
		Message:    message,
		Code:       code,
		StatusCode: status,
	}
}

func (err *RestError) Error() string {
	return fmt.Sprintf("%s(%d): %s", err.Code, err.StatusCode, err.Message)
}

func RestBadRequest(message string) *RestError {
	return NewRestError(message, "bad_request", 400)
}

func RestInternalServer(message string) *RestError {
	return NewRestError(message, "internal_error", 500)
}

func RestNotFound(message string) *RestError {
	return NewRestError(message, "not_found", 404)
}
