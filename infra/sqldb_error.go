package infra

const (
	ErrDbInternalError = iota
	ErrDbNotFound      = iota
)

type SqlDbError struct {
	Code    int
	Message string
}

func NewSqlDbError(code int, message string) *SqlDbError {
	return &SqlDbError{
		Code:    code,
		Message: message,
	}
}

func (err *SqlDbError) Error() string {
	return err.Message
}

func NewSqlDbErrInternal(message string) *SqlDbError {
	return NewSqlDbError(ErrDbInternalError, message)
}

func NewSqlDbErrNotFound(message string) *SqlDbError {
	return NewSqlDbError(ErrDbNotFound, message)
}
