package auth

type LoginRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type LoginResponse struct {
	Message string
	Token   string
}

type RegisterRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Name     string `binding:"required"`
	Email    string `binding:"required"`
}

type RegisterResponse struct {
	Message string
}
