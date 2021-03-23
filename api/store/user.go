package store

import (
	"database/sql"

	"github.com/bira37/go-rest-api/api/domain/db"
	"github.com/bira37/go-rest-api/api/domain/user"
)

type User struct{}

func NewUser() user.Store {
	return &User{}
}

func (r *User) Insert(user user.Model, connection db.Connection) (user.Model, error) {
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
		return user, db.ErrDBInternal(err.Error())
	}

	return user, err
}

func (r *User) FindByUsername(username string, connection db.Connection) (user.Model, error) {
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
		return user, db.ErrDBNotFound("User not found.")
	default:
		return user, db.ErrDBInternal(err.Error())
	}
}
