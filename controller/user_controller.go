package controller

import (
	controllerContract "bira.io/template/contract/controller"
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/dto"
	"github.com/gin-gonic/gin"
)

type userController struct {
	userService serviceContract.UserService
}

func NewUserController(userService serviceContract.UserService) controllerContract.UserController {
	return &userController{
		userService: userService,
	}
}

func (ctrl *userController) GetMe(c *gin.Context) {
	username := c.GetString("username")

	getLoggedUserRequest := dto.GetLoggedUserRequest{
		Username: username,
	}

	var getLoggedUserResponse dto.GetLoggedUserResponse

	getLoggedUserResponse, err := ctrl.userService.GetLoggedUser(getLoggedUserRequest)

	SendResponse(c, getLoggedUserResponse, err)
}
