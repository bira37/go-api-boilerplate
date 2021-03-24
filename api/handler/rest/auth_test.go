package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bira37/go-rest-api/api/domain/auth"
	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/api/mock"
	"github.com/bira37/go-rest-api/api/store"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	tmock "github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		username           string
		password           string
		expectedExists     bool
		expectedResult     user.Model
		expectedStatusCode int
		expectedErrorCode  string
		expectedStoreErr   error
		expectedToken      string
		success            bool
	}{
		{
			username:           "correct",
			password:           "correctpassword",
			expectedExists:     true,
			expectedResult:     GenUser("correct", "correctpassword"),
			expectedStatusCode: 200,
			expectedStoreErr:   nil,
			expectedToken:      "token",
			success:            true,
		},
		{
			username:           "correct",
			password:           "wrongpassword",
			expectedExists:     true,
			expectedResult:     GenUser("correct", "wrong"),
			expectedStatusCode: 400,
			expectedErrorCode:  "bad_request",
			success:            false,
		},
		{
			username:           "wrong",
			password:           "invalid",
			expectedExists:     false,
			expectedResult:     user.Model{},
			expectedStatusCode: 404,
			expectedErrorCode:  "not_found",
			expectedStoreErr:   store.ErrDBNotFound(""),
			success:            false,
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		mockDb := new(mock.MockDB)
		mockUserStore := new(mock.MockUserStore)
		jwtParser := jwt.NewJwt(Config.JwtSigningSecret)

		handler := NewAuth(mockDb, mockUserStore)

		mockUserStore.On("FindByUsername", tmock.Anything, tc.username).Return(tc.expectedResult, tc.expectedStoreErr)

		r.POST("/test-login", handler.Login)

		body, err := json.Marshal(auth.LoginRequest{
			Username: tc.username,
			Password: tc.password,
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		c.Request, err = http.NewRequest(http.MethodPost, "/test-login", bytes.NewBuffer(body))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		r.ServeHTTP(res, c.Request)

		if res.Result().StatusCode != tc.expectedStatusCode {
			t.Errorf("expected status %v, got %v", tc.expectedStatusCode, res.Result().StatusCode)
		}

		if tc.expectedStatusCode == 200 {
			var loginResponse auth.LoginResponse

			err = json.NewDecoder(res.Body).Decode(&loginResponse)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			tokenClaims, err := jwtParser.ParseToken(loginResponse.Token)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tokenClaims["sub"] != tc.username {
				t.Errorf("invalid token. expected user %v, got %v", tc.username, tokenClaims["sub"])
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
