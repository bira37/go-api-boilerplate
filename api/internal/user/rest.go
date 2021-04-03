package user

import (
	"time"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/bira37/go-rest-api/api/internal/rest"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/bira37/go-rest-api/pkg/password"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestHandler interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Me(c *gin.Context)
}

type restHandler struct {
	DB        cockroach.DB
	UserStore Store
}

func NewRestHandler(db cockroach.DB, us Store) *restHandler {
	return &restHandler{
		DB:        db,
		UserStore: us,
	}
}

func (r *restHandler) Me(ctx *gin.Context) {
	username := ctx.GetString("username")

	connection := r.DB.GetConnection()

	dbUser, err := r.UserStore.FindByUsername(connection, username)

	response := MeResponse{
		Id:       dbUser.Id,
		Username: dbUser.Username,
		Name:     dbUser.Name,
		Email:    dbUser.Email,
	}

	rest.SetResponse(ctx, response, err)
}

func (r *restHandler) Login(ctx *gin.Context) {
	var request LoginRequest
	var response LoginResponse

	if err := rest.ParseBody(ctx, &request); err != nil {
		rest.SetResponse(ctx, response, err)
		return
	}

	connection := r.DB.GetConnection()

	user, err := r.UserStore.FindByUsername(connection, request.Username)

	if err != nil {
		rest.SetResponse(ctx, response, err)
		return
	}

	if equal := password.CheckPassword(request.Password, user.PasswordHash); !equal {
		rest.SetResponse(ctx, response, errs.RestBadRequest("Invalid login or password."))
		return
	}

	jwtParser := jwt.NewJwt(config.JwtSigningSecret)

	token, err := jwtParser.GenerateToken(user.Username, 20, make(map[string]interface{}))

	if err != nil {
		rest.SetResponse(ctx, response, err)
	}

	response = LoginResponse{
		Message: "Hello " + user.Name,
		Token:   token,
	}

	rest.SetResponse(ctx, response, nil)
}

func (r *restHandler) Register(ctx *gin.Context) {
	var request RegisterRequest
	var response RegisterResponse

	if err := rest.ParseBody(ctx, &request); err != nil {
		rest.SetResponse(ctx, response, err)
		return
	}

	connection := r.DB.GetConnection()

	_, err := r.UserStore.FindByUsername(connection, request.Username)

	if err == nil {
		rest.SetResponse(ctx, response, errs.RestBadRequest("An user with the same username already exists."))
		return
	}

	passwordHash, err := password.HashPassword(request.Password)

	if err != nil {
		rest.SetResponse(ctx, response, err)
		return
	}

	newUser := Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         request.Name,
		Username:     request.Username,
		PasswordHash: passwordHash,
		Email:        request.Email,
	}

	_, err = r.UserStore.Insert(connection, newUser)

	if err != nil {
		rest.SetResponse(ctx, response, err)
		return
	}

	response = RegisterResponse{
		Message: "Registered " + request.Name,
	}

	rest.SetResponse(ctx, response, nil)
}
