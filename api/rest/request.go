package rest

import (
	"encoding/json"
	"fmt"

	"github.com/bira37/go-rest-api/api/errs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
		return errs.NewRestError(message, "bad_request", 400)
	}
	return nil
}
