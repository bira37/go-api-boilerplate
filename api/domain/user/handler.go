package user

import "github.com/gin-gonic/gin"

type RestHandler interface {
	Me(c *gin.Context)
}
