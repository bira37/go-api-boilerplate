package httprouter

import (
	controllerContract "bira.io/template/contract/controller"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(r *gin.Engine, ctrl controllerContract.AuthController) {
	public := r.Group("/auth")
	{
		public.POST("/register", ctrl.Register)
		public.POST("/login", ctrl.Login)
	}
}
