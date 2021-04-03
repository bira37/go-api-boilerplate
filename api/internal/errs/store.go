package errs

import "fmt"

type StoreError struct {
	Code    string
	Message string
}

func NewStoreError(code string, message string) *StoreError {
	return &StoreError{
		Code:    code,
		Message: message,
	}
}

func StoreInternal(message string) *StoreError {
	return NewStoreError("StoreErrorInternal", message)
}

func StoreNotFound(message string) *StoreError {
	return NewStoreError("StoreErrorNotFound", message)
}

func (err *StoreError) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

func (err *StoreError) ToRestError() *RestError {
	switch err.Code {
	case StoreNotFound("").Code:
		return RestNotFound(err.Message)
	case StoreInternal("").Code:
		return RestInternalServer(err.Message)
	default:
		return RestInternalServer("Internal error.")
	}
}
