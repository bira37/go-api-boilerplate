package rest

import (
	"encoding/json"
	"fmt"

	"github.com/bira37/go-rest-api/api/store"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Error struct {
	Code       string
	Message    string
	StatusCode int
}

func NewError(message string, code string, status int) *Error {
	return &Error{
		Message:    message,
		Code:       code,
		StatusCode: status,
	}
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s(%d): %s", err.Code, err.StatusCode, err.Message)
}

func ErrBadRequest(message string) *Error {
	return NewError(message, "bad_request", 400)
}

func ErrInternalServer(message string) *Error {
	return NewError(message, "internal_error", 500)
}

func ErrNotFound(message string) *Error {
	return NewError(message, "not_found", 404)
}

func ParseBody(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		var message string

		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			message = buildUnmarshallTypeError(e)
		case validator.ValidationErrors:
			message = buildValidationError(e)
		default:
			fmt.Printf("Unrecognized error type: %T.", e)
			message = err.Error()
		}

		return NewError(message, "bad_request", 400)
	}
	return nil
}

func SetResponse(c *gin.Context, obj interface{}, err error) {
	if err == nil {
		returnSuccess(c, obj)
	} else {
		returnError(c, err)
	}
}

func returnSuccess(c *gin.Context, obj interface{}) {
	c.JSON(200, obj)
}

func returnError(c *gin.Context, err error) {
	httpError, parseOk := err.(*Error)

	if !parseOk {
		storeErr, storeErrOk := err.(*store.Error)
		if storeErrOk {
			httpError = mapStoreError(storeErr)
		} else {
			httpError = ErrInternalServer(err.Error())
		}
	}

	c.JSON(httpError.StatusCode, map[string]interface{}{
		"message": httpError.Message,
		"code":    httpError.Code,
	})
}

func buildUnmarshallTypeError(err *json.UnmarshalTypeError) string {
	return fmt.Sprintf("'%s' should have type '%s', but has type '%s'.", err.Field, err.Type, err.Value)
}

func buildValidationError(err validator.ValidationErrors) string {
	var message string
	if len(err) > 1 {
		message += "Several errors occurred:\n"
	}
	for _, elem := range err {
		if elem.ActualTag() == "required" {
			message += fmt.Sprintf("'%s' is required.\n", elem.Field())
		} else {
			message += fmt.Sprintf("Value for '%s' is invalid.\n", elem.Field())
		}
	}
	return message
}

func mapStoreError(err *store.Error) *Error {
	switch err.Code {
	case store.ErrDBNotFound("").Code:
		return ErrNotFound(err.Message)
	case store.ErrDBInternal("").Code:
		return ErrInternalServer(err.Message)
	default:
		return ErrInternalServer("Internal error.")
	}
}
