package rest

import (
	"encoding/json"
	"fmt"

	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ParseBody(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		var message string

		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			message = formatUnmarshallTypeError(e)
		case validator.ValidationErrors:
			message = formatValidationErrors(e)
		default:
			fmt.Printf("Unrecognized error type: %T.", e)
			message = err.Error()
		}
		return errs.NewRestError(message, "bad_request", 400)
	}
	return nil
}

func formatUnmarshallTypeError(err *json.UnmarshalTypeError) string {
	return fmt.Sprintf("'%s' should have type '%s', but has type '%s'.", err.Field, err.Type, err.Value)
}

func formatValidationErrors(err validator.ValidationErrors) string {
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
