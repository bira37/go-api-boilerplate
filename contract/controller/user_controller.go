package contract

import "github.com/gin-gonic/gin"

type UserController interface {
	GetMe(c *gin.Context)
}
