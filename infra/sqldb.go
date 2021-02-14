package infra

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SqlDb interface {
	Transaction(fn func(*sqlx.Tx) error) error
	Execute(fn func(*sqlx.DB) error) error
}

type sqlDb struct {
	db *sqlx.DB
}

func NewSqlDb(databaseName string) SqlDb {
	return &sqlDb{
		db: sqlx.MustConnect("postgres", BuildSqlConnectionString("postgres", Config.SqlDbConnectionString, databaseName, "sslmode=disable")),
	}
}

func (s *sqlDb) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := s.db.Beginx()

	if err != nil {
		return NewSqlDbErrInternal(err.Error())
	}

	if err := fn(tx); err != nil {
		defer func() {
			if err := tx.Rollback(); err != nil {
				fmt.Println(err.Error())
			}
		}()
		return NewSqlDbErrInternal(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return NewSqlDbErrInternal(err.Error())
	}
	return nil
}

func (s *sqlDb) Execute(fn func(*sqlx.DB) error) error {
	return fn(s.db)
}

func BuildSqlConnectionString(driver string, uri string, dbName string, optionString string) string {
	return fmt.Sprintf("%s://%s/%s?%s", driver, uri, dbName, optionString)
}
