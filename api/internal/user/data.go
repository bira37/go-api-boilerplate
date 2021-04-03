package user

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name         string `binding:"required"`
	Username     string `binding:"required"`
	PasswordHash string `binding:"required"`
	Email        string `binding:"required"`
}

type MeResponse struct {
	Id       uuid.UUID
	Username string
	Name     string
	Email    string
}

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
