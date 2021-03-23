package server

import (
	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/handler/rest"
	"github.com/bira37/go-rest-api/api/middleware"
	"github.com/bira37/go-rest-api/api/store"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// Main Configuration
var (
	Config      = config.GetConfig()
	CockroachDB = cockroach.NewCockroachDB(Config.SQLDBConnectionString)
	JwtParser   = jwt.NewJwt(Config.JwtSigningSecret)
)

// Middlewares
var (
	AuthMiddleware = middleware.NewAuthMiddleware(JwtParser)
)

// Stores
var (
	UserStore = store.NewUser()
)

// Handlers
var (
	UserRestHandler = rest.NewUser(CockroachDB, UserStore)
	AuthRestHandler = rest.NewAuth(CockroachDB, UserStore, JwtParser)
)

// SetupServer setups middlewares, routes and handlers, returning a ready-to-start server
func SetupServer() *gin.Engine {
	router := gin.Default()

	SetupRoutes(router)

	return router
}

// SetupRoutes adds routes and middlewares to the server engine
func SetupRoutes(r *gin.Engine) {

	public := r.Group("/")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", AuthRestHandler.Login)
			auth.POST("/register", AuthRestHandler.Register)
		}
	}

	private := r.Group("/")
	{
		private.Use(AuthMiddleware)
		user := private.Group("/user")
		{
			user.GET("/me", UserRestHandler.Me)
		}
	}
}
