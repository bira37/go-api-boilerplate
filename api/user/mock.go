package user

import (
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{}
}

func (s *MockUserStore) Insert(connection cockroach.Connection, model Model) (Model, error) {
	args := s.Called(connection, model)
	return args.Get(0).(Model), args.Error(1)
}

func (s *MockUserStore) FindByUsername(connection cockroach.Connection, username string) (Model, error) {
	args := s.Called(connection, username)
	return args.Get(0).(Model), args.Error(1)
}
