package rest

import (
	"time"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/domain/auth"
	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/bira37/go-rest-api/pkg/password"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Auth struct {
	DB        cockroach.DB
	UserStore user.Store
}

var Config config.Config = config.GetConfig()

func NewAuth(db cockroach.DB, us user.Store) *Auth {
	return &Auth{
		DB:        db,
		UserStore: us,
	}
}

func (h *Auth) Login(ctx *gin.Context) {
	var request auth.LoginRequest
	var response auth.LoginResponse

	if err := ParseBody(ctx, &request); err != nil {
		SetResponse(ctx, response, err)
		return
	}

	connection := h.DB.GetConnection()

	user, err := h.UserStore.FindByUsername(connection, request.Username)

	if err != nil {
		SetResponse(ctx, response, err)
		return
	}

	if equal := password.CheckPassword(request.Password, user.PasswordHash); !equal {
		SetResponse(ctx, response, ErrBadRequest("Invalid login or password."))
		return
	}

	jwtParser := jwt.NewJwt(Config.JwtSigningSecret)

	token, err := jwtParser.GenerateToken(user.Username, 20, make(map[string]interface{}))

	if err != nil {
		SetResponse(ctx, response, err)
	}

	response = auth.LoginResponse{
		Message: "Hello " + user.Name,
		Token:   token,
	}

	SetResponse(ctx, response, nil)
}

func (h *Auth) Register(ctx *gin.Context) {
	var request auth.RegisterRequest
	var response auth.RegisterResponse

	if err := ParseBody(ctx, &request); err != nil {
		SetResponse(ctx, response, err)
		return
	}

	connection := h.DB.GetConnection()

	_, err := h.UserStore.FindByUsername(connection, request.Username)

	if err == nil {
		SetResponse(ctx, response, ErrBadRequest("An user with the same username already exists."))
		return
	}

	passwordHash, err := password.HashPassword(request.Password)

	if err != nil {
		SetResponse(ctx, response, err)
		return
	}

	newUser := user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         request.Name,
		Username:     request.Username,
		PasswordHash: passwordHash,
		Email:        request.Email,
	}

	_, err = h.UserStore.Insert(connection, newUser)

	if err != nil {
		SetResponse(ctx, response, err)
		return
	}

	response = auth.RegisterResponse{
		Message: "Registered " + request.Name,
	}

	SetResponse(ctx, response, nil)
}
