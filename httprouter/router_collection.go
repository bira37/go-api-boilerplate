package httprouter

import (
	"bira.io/template/controller"
	"github.com/gin-gonic/gin"
)

func AddRouters(r *gin.Engine, cc *controller.ControllerCollection) {
	AddAuthRouter(r, cc.AuthController)
	AddUserRouter(r, cc.UserController)
}
