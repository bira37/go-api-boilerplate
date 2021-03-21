package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	SigningSecret string
}

func (j *Jwt) GenerateToken(username string, expirationMinutes int, extraFields map[string]interface{}) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	mapClaims := jwt.MapClaims{
		"exp": time.Now().UTC().Add(time.Minute * time.Duration(expirationMinutes)).Unix(),
		"iat": time.Now().UTC(),
		"sub": username,
	}

	for k, v := range extraFields {
		mapClaims[k] = v
	}

	token.Claims = mapClaims

	tokenString, err := token.SignedString([]byte(j.SigningSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *Jwt) ParseToken(tokenString string) (map[string]interface{}, error) {
	var claims map[string]interface{}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SigningSecret), nil
	})

	if err != nil {
		return claims, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return claims, errors.New("Invalid token.")
	}

	return claims, nil
}
