package rest

import "github.com/bira37/go-rest-api/pkg/cockroach"

func ClearData(DB *cockroach.CockroachDB) {
	DB.GetConnection().MustExec("DELETE FROM users WHERE TRUE;")
}
