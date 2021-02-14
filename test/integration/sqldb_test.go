package integration

import (
	"bira.io/template/infra"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("SqlDb", func() {
	It("Should run migrations and rollback successfully", func() {
		sqldb := infra.NewSqlDb(infra.Config.SqlDbName)

		_ = sqldb.Execute(func(db *sqlx.DB) error {
			db.MustExec(`CREATE DATABASE IF NOT EXISTS migrationcheckup;`)
			return nil
		})

		connString := infra.BuildSqlConnectionString("cockroach", infra.Config.SqlDbConnectionString, "migrationcheckup", "sslmode=disable")
		infra.MigrateSqlDatabase("file://../../infra/migrations", connString, true)

		infra.MigrateSqlDatabase("file://../../infra/migrations", connString, false)

		infra.MigrateSqlDatabase("file://../../infra/migrations", connString, true)

		_ = sqldb.Execute(func(db *sqlx.DB) error {
			db.MustExec(`DROP DATABASE IF EXISTS migrationcheckup;`)
			return nil
		})
	})
})
