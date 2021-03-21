package auth

import "github.com/gin-gonic/gin"

type RestHandler interface {
	Login(*gin.Context)
	Register(*gin.Context)
}
