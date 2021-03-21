package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateTokenSuccessfully(t *testing.T) {
	jwt := Jwt{SigningSecret: uuid.New().String()}

	token, err := jwt.GenerateToken("username", 20, map[string]interface{}{"field": "value"})

	if err != nil {
		t.Errorf("expected nil err, found '%v'", err)
	}

	if len(token) == 0 {
		t.Errorf("token has len zero, expected greater or equal than one")
	}
}

func TestParseTokenSuccessfully(t *testing.T) {
	jwt := Jwt{SigningSecret: uuid.New().String()}

	tokenString, _ := jwt.GenerateToken("username", 20, map[string]interface{}{"field": "value"})
	token, err := jwt.ParseToken(tokenString)

	if err != nil {
		t.Errorf("expected nil err, found %v", err)
	}

	if token["field"] != "value" {
		t.Errorf("token['field']: expected 'value', found '%v'", token["field"])
	}
}

func TestParseTokenExpired(t *testing.T) {
	jwt := Jwt{SigningSecret: uuid.New().String()}

	tokenString, _ := jwt.GenerateToken("username", 0, map[string]interface{}{"exp": time.Now().UTC().Add(-time.Second * 10).Unix()})

	_, err := jwt.ParseToken(tokenString)

	if err == nil {
		t.Errorf("expected err, found '%v'", err)
	}
}
