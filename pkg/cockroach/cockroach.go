package cockroach

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CockroachDB struct {
	db *sqlx.DB
}

func NewCockroachDB(connectionString string) *CockroachDB {
	return &CockroachDB{
		db: sqlx.MustConnect("postgres", connectionString),
	}
}

func (s *CockroachDB) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := s.db.Beginx()

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		defer func() {
			if err := tx.Rollback(); err != nil {
				fmt.Println(err.Error())
			}
		}()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *CockroachDB) GetConnection() *sqlx.DB {
	return s.db
}
