package rest

import (
	"github.com/bira37/go-rest-api/api/domain/db"
	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/gin-gonic/gin"
)

type User struct {
	DB        db.DB
	UserStore user.Store
}

func NewUser(db db.DB, us user.Store) *User {
	return &User{
		DB:        db,
		UserStore: us,
	}
}

func (h *User) Me(ctx *gin.Context) {
	username := ctx.GetString("username")

	connection := h.DB.GetConnection()

	dbUser, err := h.UserStore.FindByUsername(username, connection)

	response := user.MeResponse{
		Id:       dbUser.Id,
		Username: dbUser.Username,
		Name:     dbUser.Name,
		Email:    dbUser.Email,
	}

	SetResponse(ctx, response, err)
}
