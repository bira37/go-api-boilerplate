package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/api/mock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	tmock "github.com/stretchr/testify/mock"
)

func GenUser(username string) user.Model {
	faker := gofakeit.NewCrypto()

	return user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         faker.Name(),
		Username:     username,
		PasswordHash: faker.Password(true, true, true, false, false, 20),
		Email:        faker.Email(),
	}
}

func TestMe(t *testing.T) {
	res := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(res)

	mockDb := new(mock.MockDB)
	mockUserStore := new(mock.MockUserStore)

	handler := NewUser(mockDb, mockUserStore)

	username := "test"

	userModel := GenUser(username)

	mockUserStore.On("FindByUsername", username, tmock.Anything).Return(userModel, nil)

	r.Use(func(c *gin.Context) {
		c.Set("username", username)
		c.Next()
	})

	r.GET("/test-me", handler.Me)

	var err error

	c.Request, err = http.NewRequest(http.MethodGet, "/test-me", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	r.ServeHTTP(res, c.Request)

	if res.Result().StatusCode != 200 {
		t.Errorf("error on handler: %v", res.Body.String())
	}

	var meResponse user.MeResponse

	err = json.NewDecoder(res.Body).Decode(&meResponse)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if meResponse.Email != userModel.Email ||
		meResponse.Id != userModel.Id ||
		meResponse.Name != userModel.Name ||
		meResponse.Username != userModel.Username {
		t.Errorf("error: expected equivalent to '%v', got '%v'", userModel, meResponse)
	}

	mockUserStore.AssertExpectations(t)
}
