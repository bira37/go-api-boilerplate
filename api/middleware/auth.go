package middleware

import (
	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

var Config config.Config = config.GetConfig()

func NewAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("X-Access-Token")

		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(401, map[string]string{
				"code":    "unauthorized",
				"message": "No token provided.",
			})
			return
		}

		jwtParser := jwt.NewJwt(Config.JwtSigningSecret)

		claims, err := jwtParser.ParseToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(401, map[string]string{
				"code":    "unauthorized",
				"message": "Invalid token.",
			})
			return
		}

		c.Set("username", claims["sub"])
		c.Next()
	}
}
