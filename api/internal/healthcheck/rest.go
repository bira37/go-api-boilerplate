package healthcheck

import (
	"github.com/bira37/go-rest-api/api/internal/rest"
	"github.com/gin-gonic/gin"
)

type RestHandler interface {
	Health(*gin.Context)
}

type restHandler struct{}

func NewRestHandler() *restHandler {
	return &restHandler{}
}

func (r *restHandler) Health(c *gin.Context) {
	rest.SetResponse(c, struct{ Message string }{Message: "OK"}, nil)
}
