package mock

import (
	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (s *MockUserStore) Insert(connection cockroach.Connection, model user.Model) (user.Model, error) {
	args := s.Called(model, connection)
	return args.Get(0).(user.Model), args.Error(1)
}

func (s *MockUserStore) FindByUsername(connection cockroach.Connection, username string) (user.Model, error) {
	args := s.Called(username, connection)
	return args.Get(0).(user.Model), args.Error(1)
}
