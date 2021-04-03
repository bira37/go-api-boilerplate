package rest

import (
	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string
	Code    string
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
	httpError, parseOk := err.(*errs.RestError)

	if !parseOk {
		storeErr, storeErrOk := err.(*errs.StoreError)
		if storeErrOk {
			httpError = storeErr.ToRestError()
		} else {
			httpError = errs.RestInternalServer(err.Error())
		}
	}

	c.JSON(httpError.StatusCode, ErrorResponse{Message: httpError.Message, Code: httpError.Code})
}
