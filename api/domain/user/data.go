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
