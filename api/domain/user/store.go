package user

import (
	"time"

	"github.com/bira37/go-rest-api/api/domain/db"
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
	FindByUsername(username string, connection db.Connection) (Model, error)
	Insert(user Model, connection db.Connection) (Model, error)
}
