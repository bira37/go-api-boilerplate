package repository

import (
	"fmt"

	repositoryContract "bira.io/template/contract/repository"
	"bira.io/template/infra"
	"bira.io/template/model"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	sqlDb infra.SqlDb
}

func NewUserRepository(infraCollection infra.InfraCollection) repositoryContract.UserRepository {
	return &userRepository{
		sqlDb: infraCollection.GetSqlDb(),
	}
}

func (r *userRepository) InsertUser(userCreate model.UserCreate) (model.User, error) {
	var user model.User
	err := r.sqlDb.Execute(func(db *sqlx.DB) error {
		result, err := db.NamedQuery(`
			INSERT INTO users (name, username, password_hash, email)
			VALUES (
				:name,
				:username,
				:password_hash,
				:email
			)
			RETURNING *;
		`, userCreate)

		if err != nil {
			return infra.NewSqlDbErrInternal(err.Error())
		}

		if ok := result.Next(); ok {
			err = result.StructScan(&user)
			if err != nil {
				return infra.NewSqlDbErrInternal(err.Error())
			}
			return nil
		}

		return infra.NewSqlDbErrInternal(err.Error())
	})

	return user, err
}

func (r *userRepository) FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.sqlDb.Execute(func(db *sqlx.DB) error {
		result, err := db.Queryx(`
			SELECT *
			FROM users
			WHERE username = $1
			LIMIT 1;
		`, username)

		if err != nil {
			fmt.Println(err.Error())
			return infra.NewSqlDbErrInternal(err.Error())
		}

		if ok := result.Next(); ok {
			err = result.StructScan(&user)
			if err != nil {
				return infra.NewSqlDbErrInternal(err.Error())
			}
			return nil
		}

		return infra.NewSqlDbErrNotFound("User not found.")
	})

	return user, err
}
