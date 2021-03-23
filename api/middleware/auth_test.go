package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	res := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(res)

	r.Use(NewAuthMiddleware())

	r.GET("/test", func(c *gin.Context) {
		user, exists := c.Get("username")

		if !exists {
			t.Error("username should exist")
			c.AbortWithStatus(500)
			return
		}

		if user != "testuser" {
			t.Errorf("expected testuser, found %s", user)
			c.AbortWithStatus(500)
			return
		}
		c.Status(200)
	})

	var err error

	c.Request, err = http.NewRequest(http.MethodGet, "/test", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	jwtParser := jwt.Jwt{SigningSecret: Config.JwtSigningSecret}

	token, err := jwtParser.GenerateToken("testuser", 60, make(map[string]interface{}))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	c.Request.Header.Set("X-Access-Token", token)

	r.ServeHTTP(res, c.Request)

	if res.Result().StatusCode != 200 {
		t.Errorf("error on middleware: %v", res.Body.String())
	}
}
