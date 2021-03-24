package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/api/mock"
	"github.com/bira37/go-rest-api/api/store"
	"github.com/bira37/go-rest-api/pkg/password"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	tmock "github.com/stretchr/testify/mock"
)

func GenUser(username string, passwd string) user.Model {
	faker := gofakeit.NewCrypto()
	hash, err := password.HashPassword(passwd)

	if err != nil {
		panic(err)
	}

	return user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         faker.Name(),
		Username:     username,
		PasswordHash: hash,
		Email:        faker.Email(),
	}
}

func TestMe(t *testing.T) {
	tests := []struct {
		username           string
		expectedExists     bool
		expectedResult     user.Model
		expectedStoreErr   error
		expectedErrorCode  string
		expectedStatusCode int
	}{
		{
			username:           "success",
			expectedExists:     true,
			expectedResult:     GenUser("success", "password"),
			expectedStoreErr:   nil,
			expectedStatusCode: 200,
		},
		{
			username:           "fail",
			expectedExists:     false,
			expectedResult:     user.Model{},
			expectedStoreErr:   store.ErrDBNotFound(""),
			expectedErrorCode:  "not_found",
			expectedStatusCode: 404,
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		mockDb := new(mock.MockDB)
		mockUserStore := new(mock.MockUserStore)

		handler := NewUser(mockDb, mockUserStore)

		mockUserStore.On("FindByUsername", tmock.Anything, tc.username).Return(tc.expectedResult, tc.expectedStoreErr)

		r.Use(func(c *gin.Context) {
			c.Set("username", tc.username)
			c.Next()
		})

		r.GET("/test-me", handler.Me)

		var err error

		c.Request, err = http.NewRequest(http.MethodGet, "/test-me", nil)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		r.ServeHTTP(res, c.Request)

		if res.Result().StatusCode != tc.expectedStatusCode {
			t.Errorf("expected status %v, got %v", tc.expectedStatusCode, res.Result().StatusCode)
		}

		if tc.expectedStatusCode == 200 {
			var meResponse user.MeResponse

			err = json.NewDecoder(res.Body).Decode(&meResponse)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if meResponse.Email != tc.expectedResult.Email ||
				meResponse.Id != tc.expectedResult.Id ||
				meResponse.Name != tc.expectedResult.Name ||
				meResponse.Username != tc.expectedResult.Username {
				t.Errorf("error: expected equivalent to '%v', got '%v'", tc.expectedResult, meResponse)
			}
		} else {
			var errorMsg ErrorResponse

			err = json.NewDecoder(res.Body).Decode(&errorMsg)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if errorMsg.Code != tc.expectedErrorCode {
				t.Errorf("expected code %v, got %v", tc.expectedErrorCode, errorMsg.Code)
			}
		}

		mockUserStore.AssertExpectations(t)
	}
}
