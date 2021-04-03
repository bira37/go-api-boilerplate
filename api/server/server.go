package server

import (
	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/internal/healthcheck"
	"github.com/bira37/go-rest-api/api/internal/middleware"
	"github.com/bira37/go-rest-api/api/internal/user"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/gin-gonic/gin"
)

// Main Configuration
var (
	CockroachDB = cockroach.NewCockroachDB(config.SQLDBConnectionString)
)

// Middlewares
var (
	AuthMiddleware = middleware.NewAuthMiddleware()
)

// Stores
var (
	UserStore = user.NewStore()
)

// Handlers
var (
	UserRestHandler   = user.NewRestHandler(CockroachDB, UserStore)
	HealthRestHandler = healthcheck.NewRestHandler()
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
		health := public.Group("/health")
		{
			health.GET("", HealthRestHandler.Health)
		}

		auth := public.Group("/auth")
		{
			auth.POST("/login", UserRestHandler.Login)
			auth.POST("/register", UserRestHandler.Register)
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
