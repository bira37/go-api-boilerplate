package middleware

import (
	"bira.io/template/infra"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(infra.Config.JwtSigningString), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(500, map[string]string{
				"code":    "internal_error",
				"message": "Error parsing token.",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok && !token.Valid {
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
