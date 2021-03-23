package user

import (
	"time"

	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/google/uuid"
)

type Model struct {
	Id           uuid.UUID `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Name         string    `db:"name"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Email        string    `db:"email"`
}

type Store interface {
	FindByUsername(connection cockroach.Connection, username string) (Model, error)
	Insert(connection cockroach.Connection, user Model) (Model, error)
}
