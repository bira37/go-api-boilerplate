package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/bira37/go-rest-api/api/internal/rest"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/bira37/go-rest-api/pkg/password"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func GenUser(username string, passwd string) Model {
	faker := gofakeit.NewCrypto()
	hash, err := password.HashPassword(passwd)

	if err != nil {
		panic(err)
	}

	return Model{
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
		expectedResult     Model
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
			expectedResult:     Model{},
			expectedStoreErr:   errs.StoreNotFound(""),
			expectedErrorCode:  "not_found",
			expectedStatusCode: 404,
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		mockDb := cockroach.NewMockDB()
		mockUserStore := NewMockUserStore()

		handler := NewRestHandler(mockDb, mockUserStore)

		mockUserStore.On("FindByUsername", mock.Anything, tc.username).Return(tc.expectedResult, tc.expectedStoreErr)

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
			var meResponse MeResponse

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
			var errorMsg rest.ErrorResponse

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

func TestLogin(t *testing.T) {
	tests := []struct {
		username           string
		password           string
		expectedExists     bool
		expectedResult     Model
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
			expectedResult:     Model{},
			expectedStatusCode: 404,
			expectedErrorCode:  "not_found",
			expectedStoreErr:   errs.StoreNotFound(""),
			success:            false,
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		mockDb := cockroach.NewMockDB()
		mockUserStore := NewMockUserStore()
		jwtParser := jwt.NewJwt(config.JwtSigningSecret)

		handler := NewRestHandler(mockDb, mockUserStore)

		mockUserStore.On("FindByUsername", mock.Anything, tc.username).Return(tc.expectedResult, tc.expectedStoreErr)

		r.POST("/test-login", handler.Login)

		body, err := json.Marshal(LoginRequest{
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
			var loginResponse LoginResponse

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
			var errorMsg rest.ErrorResponse

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

func TestRegister(t *testing.T) {
	tests := []struct {
		request            RegisterRequest
		expectedStatusCode int
		expectedStoreErr   error
		expectedResponse   RegisterResponse
		success            bool
		expectedError      rest.ErrorResponse
	}{
		{
			request:            RegisterRequest{Username: "correct", Password: "password", Name: "Test", Email: "test@example.com"},
			expectedStatusCode: 200,
			expectedStoreErr:   errs.StoreNotFound(""),
			success:            true,
			expectedResponse:   RegisterResponse{Message: "Registered Test"},
		},
		{
			request:            RegisterRequest{Username: "correct", Password: "password", Name: "Test", Email: "test@example.com"},
			expectedStatusCode: 400,
			success:            false,
			expectedStoreErr:   nil,
			expectedError:      rest.ErrorResponse{Message: "An user with the same username already exists.", Code: "bad_request"},
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		mockDb := cockroach.NewMockDB()
		mockUserStore := NewMockUserStore()

		handler := NewRestHandler(mockDb, mockUserStore)

		mockUserStore.On("FindByUsername", mock.Anything, tc.request.Username).Return(GenUser(tc.request.Username, tc.request.Password), tc.expectedStoreErr)

		if tc.success {
			mockUserStore.On("Insert", mock.Anything, mock.Anything).Return(Model{}, nil)
		}

		r.POST("/test-register", handler.Register)

		body, err := json.Marshal(tc.request)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		c.Request, err = http.NewRequest(http.MethodPost, "/test-register", bytes.NewBuffer(body))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		r.ServeHTTP(res, c.Request)

		if res.Result().StatusCode != tc.expectedStatusCode {
			t.Errorf("expected status %v, got %v", tc.expectedStatusCode, res.Result().StatusCode)
		}

		if tc.expectedStatusCode == 200 {
			var registerResponse RegisterResponse

			err = json.NewDecoder(res.Body).Decode(&registerResponse)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if registerResponse.Message != tc.expectedResponse.Message {
				t.Errorf("unexpected message. expected '%v', but found '%v'", tc.expectedResponse.Message, registerResponse.Message)
			}
		} else {
			var errorMsg rest.ErrorResponse

			err = json.NewDecoder(res.Body).Decode(&errorMsg)

			if err != nil {
				t.Errorf("unexpected error: '%v'", err)
			}

			if errorMsg.Code != tc.expectedError.Code {
				t.Errorf("expected code '%v', got '%v'", tc.expectedError.Code, errorMsg.Code)
			}

			if errorMsg.Message != tc.expectedError.Message {
				t.Errorf("expected message '%v', got '%v'", tc.expectedError.Message, errorMsg.Message)
			}
		}

		mockUserStore.AssertExpectations(t)
	}
}
