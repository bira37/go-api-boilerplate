package httprouter

import (
	controllerContract "bira.io/template/contract/controller"
	"bira.io/template/middleware"
	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.Engine, ctrl controllerContract.UserController) {
	private := r.Group("/user")
	{
		private.Use(middleware.NewAuthMiddleware())
		private.GET("/me", ctrl.GetMe)
	}
}
