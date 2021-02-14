package dto

import "github.com/google/uuid"

type GetLoggedUserRequest struct {
	Username string `binding:"required"`
}

type GetLoggedUserResponse struct {
	Username string    `binding:"required"`
	Name     string    `binding:"required"`
	Email    string    `binding:"required"`
	Id       uuid.UUID `binding:"required"`
}
