package store

import (
	"fmt"
	"os"
	"testing"

	"github.com/bira37/go-rest-api/pkg/cockroach"
)

var Connection *cockroach.CockroachDB

func TestMain(m *testing.M) {
	Connection = cockroach.NewCockroachDB("postgres://root@localhost:26257?sslmode=disable")

	err := Connection.GetConnection().Ping()

	if err != nil {
		panic(fmt.Sprintf("Could not establish connection: %v", err))
	}

	code := m.Run()

	os.Exit(code)
}
