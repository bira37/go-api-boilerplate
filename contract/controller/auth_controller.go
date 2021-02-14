package contract

import "github.com/gin-gonic/gin"

type AuthController interface {
	Login(*gin.Context)
	Register(*gin.Context)
}
