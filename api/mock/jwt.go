package mock

import "github.com/stretchr/testify/mock"

type MockJwt struct {
	mock.Mock
}

func (j *MockJwt) GenerateToken(username string, expirationMinutes int, extraFields map[string]interface{}) (string, error) {
	args := j.Called(username, expirationMinutes, extraFields)
	return args.String(0), args.Error(1)
}

func (j *MockJwt) ParseToken(tokenString string) (map[string]interface{}, error) {
	args := j.Called(tokenString)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
