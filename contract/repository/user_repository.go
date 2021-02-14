package contract

import (
	"bira.io/template/model"
)

type UserRepository interface {
	FindUserByUsername(username string) (model.User, error)
	InsertUser(userCreate model.UserCreate) (model.User, error)
}
