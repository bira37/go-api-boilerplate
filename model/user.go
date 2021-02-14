package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `binding:"required" db:"id"`
	CreatedAt    time.Time `binding:"required" db:"created_at"`
	UpdatedAt    time.Time `binding:"required" db:"updated_at"`
	Name         string    `binding:"required" db:"name"`
	Username     string    `binding:"required" db:"username"`
	PasswordHash string    `binding:"required" db:"password_hash"`
	Email        string    `binding:"required" db:"email"`
}

type UserCreate struct {
	Name         string `binding:"required" db:"name"`
	Username     string `binding:"required" db:"username"`
	PasswordHash string `binding:"required" db:"password_hash"`
	Email        string `binding:"required" db:"email"`
}
