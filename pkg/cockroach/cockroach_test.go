package cockroach

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var Connection *CockroachDB

type TestItem struct {
	Id      uuid.UUID
	Content string
}

func GenTestItem() TestItem {
	return TestItem{
		Id:      uuid.New(),
		Content: "test",
	}
}

func TestMain(m *testing.M) {
	rootConnection := NewCockroachDB("postgres://root@localhost:26257?sslmode=disable")
	defer rootConnection.db.Close()

	err := rootConnection.db.Ping()

	if err != nil {
		panic(fmt.Sprintf("Could not establish connection for root: %v", err))
	}

	rootConnection.db.MustExec(`CREATE DATABASE IF NOT EXISTS cockroachtest;`)

	Connection = NewCockroachDB("postgres://root@localhost:26257/cockroachtest?sslmode=disable")
	defer Connection.db.Close()

	err = Connection.db.Ping()

	if err != nil {
		panic(fmt.Sprintf("Could not establish connection: %v", err))
	}

	Connection.db.MustExec(`CREATE TABLE IF NOT EXISTS test_table (id UUID PRIMARY KEY NOT NULL, content TEXT NOT NULL);`)

	code := m.Run()

	rootConnection.db.MustExec(`DROP DATABASE IF EXISTS cockroachtest;`)

	os.Exit(code)
}

func TestCockroachDBStatement(t *testing.T) {
	connection := Connection.GetConnection()

	values := GenTestItem()

	res, err := connection.NamedExec(`INSERT INTO test_table VALUES (:id, :content)`, values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	count, err := res.RowsAffected()

	if count != 1 || err != nil {
		t.Errorf("expected to affect one row, affected %v\nerr: %v", count, err)
	}

	res = connection.MustExec(`DELETE FROM test_table WHERE id = $1`, values.Id)

	count, err = res.RowsAffected()

	if count != 1 || err != nil {
		t.Errorf("expected to delete one row, deleted %v\nerr: %v", count, err)
	}
}

func TestCockroachDBTransaction(t *testing.T) {
	values := []TestItem{GenTestItem(), GenTestItem()}

	_ = Connection.Transaction(func(tx *sqlx.Tx) error {
		for _, value := range values {
			tx.MustExec(`INSERT INTO test_table VALUES ($1, $2)`, value.Id, value.Content)
		}
		return nil
	})

	connection := Connection.GetConnection()

	for _, value := range values {
		var item TestItem
		err := connection.Get(&item, `SELECT * FROM test_table WHERE id = $1 LIMIT 1`, value.Id)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(item, value) {
			t.Errorf("expected %v, found %v", value, item)
		}
	}

	_ = Connection.Transaction(func(tx *sqlx.Tx) error {
		for _, value := range values {
			res := tx.MustExec(`DELETE FROM test_table WHERE id = $1`, value.Id)

			rowsAffected, err := res.RowsAffected()

			if rowsAffected != 1 || err != nil {
				t.Errorf("expected to affect one row, affected %v\nerr: %v", rowsAffected, err)
			}
		}
		return nil
	})
}

func TestCockroachDBTransactionRollback(t *testing.T) {
	value := GenTestItem()

	err := Connection.Transaction(func(tx *sqlx.Tx) error {
		tx.MustExec(`INSERT INTO test_table VALUES ($1, $2)`, value.Id, value.Content)
		return fmt.Errorf("forced error")
	})

	if err == nil {
		t.Errorf("expected transaction to return error and rollback")
	}

	var item TestItem

	err = Connection.GetConnection().Get(&item, `SELECT * FROM test_table WHERE id = $1`, value.Id)

	if err == nil {
		t.Error("expected query returns no rows and return error")
	}
}
