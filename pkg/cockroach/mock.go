package cockroach

import (
	"github.com/jmoiron/sqlx"
)

type MockDB struct{}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (db *MockDB) Transaction(fn func(*sqlx.Tx) error) error {
	tx := &sqlx.Tx{}
	if err := fn(tx); err != nil {
		return err
	}
	return nil
}

func (db *MockDB) GetConnection() *sqlx.DB {
	return &sqlx.DB{}
}
