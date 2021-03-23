package store

import (
	"database/sql"

	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/pkg/cockroach"
)

type User struct{}

func NewUser() user.Store {
	return &User{}
}

func (r *User) Insert(connection cockroach.Connection, user user.Model) (user.Model, error) {
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
		return user, ErrDBInternal(err.Error())
	}

	return user, err
}

func (r *User) FindByUsername(connection cockroach.Connection, username string) (user.Model, error) {
	var user user.Model

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
		return user, ErrDBNotFound("User not found.")
	default:
		return user, ErrDBInternal(err.Error())
	}
}
