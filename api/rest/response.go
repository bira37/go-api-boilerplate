package rest

import (
	"encoding/json"
	"fmt"

	"github.com/bira37/go-rest-api/api/errs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
	httpError, parseOk := err.(*errs.RestError)

	if !parseOk {
		storeErr, storeErrOk := err.(*errs.StoreError)
		if storeErrOk {
			httpError = mapStoreError(storeErr)
		} else {
			httpError = errs.RestInternalServer(err.Error())
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

func mapStoreError(err *errs.StoreError) *errs.RestError {
	switch err.Code {
	case errs.StoreNotFound("").Code:
		return errs.RestNotFound(err.Message)
	case errs.StoreInternal("").Code:
		return errs.RestInternalServer(err.Message)
	default:
		return errs.RestInternalServer("Internal error.")
	}
}
