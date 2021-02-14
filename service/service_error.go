package service

type HttpError struct {
	Code       string
	StatusCode int
	Message    string
}

func NewHttpError(message string, code string, statusCode int) *HttpError {
	return &HttpError{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

func (err *HttpError) Error() string {
	return err.Message
}

func NewHttpErrBadRequest(message string) *HttpError {
	return NewHttpError(message, "bad_request", 400)
}

func NewHttpErrInternalServer(message string) *HttpError {
	return NewHttpError(message, "internal_error", 500)
}

func NewHttpErrNotFound(message string) *HttpError {
	return NewHttpError(message, "not_found", 404)
}
