package controller

import (
	controllerContract "bira.io/template/contract/controller"
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/dto"
	"github.com/gin-gonic/gin"
)

type authController struct {
	authService serviceContract.AuthService
}

func NewAuthController(authService serviceContract.AuthService) controllerContract.AuthController {
	return &authController{
		authService: authService,
	}
}

func (ctrl *authController) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	var loginResponse dto.LoginResponse

	if !ParseBody(c, &loginRequest) {
		return
	}

	loginResponse, err := ctrl.authService.Login(loginRequest)

	SendResponse(c, loginResponse, err)
}

func (ctrl *authController) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	var registerResponse dto.RegisterResponse

	if !ParseBody(c, &registerRequest) {
		return
	}

	registerResponse, err := ctrl.authService.Register(registerRequest)

	SendResponse(c, registerResponse, err)
}
