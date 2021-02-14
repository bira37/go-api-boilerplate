package dto

type LoginRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type LoginResponse struct {
	Message string `binding:"required"`
	Token   string `binding:"required"`
}

type RegisterRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Name     string `binding:"required"`
	Email    string `binding:"required"`
}

type RegisterResponse struct {
	Message string `binding:"required"`
}
