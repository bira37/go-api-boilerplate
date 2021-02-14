package infra

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateSqlDatabase(migrationsPath string, connString string, upMigration bool) {
	if upMigration {
		fmt.Println("Running migrations...")
	} else {
		fmt.Println("Running rollback...")
	}

	migrator, err := migrate.New(migrationsPath, connString)

	if err != nil {
		log.Fatal(err)
	}

	var step int

	if upMigration {
		step = 1
	} else {
		step = -1
	}

	for {
		version, dirty, versionErr := migrator.Version()

		fmt.Printf("Migrating database at version %d...\n", version)

		if dirty {
			log.Fatalln("Database is dirty, fix it mannualy.")
		}

		if versionErr != nil && !strings.Contains(versionErr.Error(), "no migration") {
			log.Fatalln(versionErr.Error())
		}

		if err := migrator.Steps(step); err != nil {
			if strings.Contains(err.Error(), "file does not exist") {
				break
			}

			fmt.Println("Migration failed: " + err.Error())
			fmt.Println("Trying to force down...")

			if versionErr != nil && strings.Contains(versionErr.Error(), "no migration") {
				log.Fatalln("Cannot force down the first migration. Solve it mannualy.")
			}

			err = migrator.Force(int(version))

			if err != nil {
				log.Fatalln("Error while forcing down: " + err.Error() + "\nHandle the error mannualy.")
			}

			log.Fatalln("Forced down successfully. Fix migration and try again.")
		}

		fmt.Println("Migration succeeded.")
	}

	fmt.Println("Finished migrations.")
}
