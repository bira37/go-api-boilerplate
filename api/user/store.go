package user

import (
	"database/sql"
	"time"

	"github.com/bira37/go-rest-api/api/errs"
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

type store struct{}

func NewStore() *store {
	return &store{}
}

func (s *store) Insert(connection cockroach.Connection, user Model) (Model, error) {
	_, err := connection.NamedExec(`
			INSERT INTO users (id, name, username, password_hash, email, created_at, updated_at)
			VALUES (
				:id,
				:name,
				:username,
				:password_hash,
				:email,
				:created_at,
				:updated_at
			);
		`, user)

	if err != nil {
		return user, errs.StoreInternal(err.Error())
	}

	return user, err
}

func (s *store) FindByUsername(connection cockroach.Connection, username string) (Model, error) {
	var user Model

	err := connection.Get(&user, `
			SELECT *
			FROM users
			WHERE username = $1
			LIMIT 1;
		`, username)

	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return user, errs.StoreNotFound("User not found.")
	default:
		return user, errs.StoreInternal(err.Error())
	}
}
