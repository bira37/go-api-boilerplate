package cockroach

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type DB interface {
	Transaction(fn func(*sqlx.Tx) error) error
	GetConnection() *sqlx.DB
}

type Connection interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type CockroachDB struct {
	db *sqlx.DB
}

func NewCockroachDB(connectionString string) *CockroachDB {
	return &CockroachDB{
		db: sqlx.MustConnect("postgres", connectionString),
	}
}

func (s *CockroachDB) Transaction(fn func(*sqlx.Tx) error) error {
	retries := 0
	rand.Seed(time.Now().UTC().UnixNano())

	// retry at most 50 times
	for {
		retries++

		if retries > 50 {
			fmt.Println("Transaction maximum retries reached.")
			return errors.New("Transaction maximum retries reached.")
		}

		// begin transaction
		tx, err := s.db.Beginx()

		if err != nil {
			return err
		}

		// if gets error on execution, just stop and rollback
		if err := fn(tx); err != nil {
			defer func() {
				if err := tx.Rollback(); err != nil {
					fmt.Println(err.Error())
				}
			}()
			return err
		}

		// try to commit
		err = tx.Commit()

		// check if error code was 40001 (serialization error) and retry if necessary
		if err != nil {
			if pgerr, ok := err.(*pq.Error); ok {
				if pgerr.Code == "40001" {
					timems := int(math.Pow(2, float64(retries/5))) + rand.Intn(100)
					time.Sleep(time.Duration(timems) * time.Millisecond)
					continue
				}
			}
			return err
		}
		return nil
	}
}

func (s *CockroachDB) GetConnection() *sqlx.DB {
	return s.db
}
