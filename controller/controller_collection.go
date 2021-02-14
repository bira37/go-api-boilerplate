package controller

import (
	"encoding/json"
	"fmt"

	controllerContract "bira.io/template/contract/controller"
	"bira.io/template/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ControllerCollection struct {
	AuthController controllerContract.AuthController
	UserController controllerContract.UserController
}

func NewControllerCollection(sc *service.ServiceCollection) *ControllerCollection {
	return &ControllerCollection{
		AuthController: NewAuthController(sc.AuthService),
		UserController: NewUserController(sc.UserService),
	}
}

func ParseBody(c *gin.Context, obj interface{}) bool {
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

		c.JSON(400, gin.H{
			"message": message,
			"code":    "bad_request",
		})

		return false
	}
	return true
}

func SendResponse(c *gin.Context, obj interface{}, err error) {
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
	httpError, parseOk := err.(*service.HttpError)

	if !parseOk {
		httpError = service.NewHttpErrInternalServer(err.Error())
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
