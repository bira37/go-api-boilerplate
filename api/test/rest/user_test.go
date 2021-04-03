package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/internal/user"
	"github.com/bira37/go-rest-api/pkg/jwt"
	"github.com/bira37/go-rest-api/pkg/password"
	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	body, err := json.Marshal(user.RegisterRequest{Username: "register", Password: "test", Name: "Test", Email: "register@test.com"})

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/auth/register", Server.URL), "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected 200, got %v", res.StatusCode)
	}

	var response user.RegisterResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	if response.Message != "Registered Test" {
		t.Errorf("expected 'Registered Test' message, got '%v'", response.Message)
	}
}

func TestLogin(t *testing.T) {
	userStore := user.NewStore()
	connection := DB.GetConnection()

	hash, err := password.HashPassword("test")

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	userStore.Insert(connection, user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         "test",
		Username:     "login",
		PasswordHash: hash,
		Email:        "login@example.com",
	})

	body, err := json.Marshal(user.LoginRequest{Username: "login", Password: "test"})

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/auth/login", Server.URL), "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected 200, got %v", res.StatusCode)
	}

	var response user.LoginResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	if response.Message != "Hello, test" {
		t.Errorf("expected message 'Hello, test', but found '%v'", response.Message)
	}
}

func TestMe(t *testing.T) {
	userStore := user.NewStore()
	connection := DB.GetConnection()

	hash, err := password.HashPassword("test")

	jwtParser := jwt.NewJwt(config.JwtSigningSecret)

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	userStore.Insert(connection, user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         "test",
		Username:     "me",
		PasswordHash: hash,
		Email:        "login@example.com",
	})

	token, err := jwtParser.GenerateToken("me", 15, map[string]interface{}{})

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user/me", Server.URL), nil)

	req.Header.Set("X-Access-Token", token)

	client := &http.Client{}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		t.Errorf("expected 200, got %v", res.StatusCode)
	}

	var response user.MeResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		t.Errorf("unexpected err: '%v'", err)
	}

	if response.Username != "me" {
		t.Errorf("expected username 'me', but found '%v'", response.Username)
	}
}
