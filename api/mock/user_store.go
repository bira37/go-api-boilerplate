package mock

import (
	"github.com/bira37/go-rest-api/api/domain/db"
	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (s *MockUserStore) Insert(model user.Model, connection db.Connection) (user.Model, error) {
	args := s.Called(model, connection)
	return args.Get(0).(user.Model), args.Error(1)
}

func (s *MockUserStore) FindByUsername(username string, connection db.Connection) (user.Model, error) {
	args := s.Called(username, connection)
	return args.Get(0).(user.Model), args.Error(1)
}
