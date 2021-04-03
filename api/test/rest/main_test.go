package rest

import (
	"net/http/httptest"
	"testing"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/server"
	"github.com/bira37/go-rest-api/pkg/cockroach"
)

var Server *httptest.Server
var DB *cockroach.CockroachDB

func TestMain(m *testing.M) {
	DB = cockroach.NewCockroachDB(config.SQLDBConnectionString)
	ClearData()
	Server = httptest.NewServer(server.SetupRestServer())
	defer Server.Close()
	m.Run()
}

func ClearData() {
	DB.GetConnection().MustExec("DELETE FROM users WHERE TRUE;")
}
