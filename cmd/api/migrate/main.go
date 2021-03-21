package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	os.Setenv("GOOSE_DRIVER", "postgres")
	os.Setenv("GOOSE_DBSTRING", "postgres://root@localhost:26257?sslmode=disable")

	args := os.Args[1:]

	if len(args) < 1 {
		panic(`Missing argument. Expected <create|up> <migration_name>`)
	}

	switch args[0] {
	case "create":
		if len(args) < 2 {
			panic(`Expected migration name. "create <migration_name>"`)
		}
		createMigration(args[1])
	case "up":
		migrateUp()
	default:
		panic("Command not found: expected <create|up>")
	}
}

func migrateUp() {
	cmd := exec.Command("goose", "-dir", "api/migrations", "up")

	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	cmd.Stderr = writer
	cmd.Stdout = writer

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	fmt.Println(buffer.String())
}

func createMigration(migrationName string) {
	cmd := exec.Command("goose", "-dir", "api/migrations", "create", migrationName, "sql")

	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	cmd.Stderr = writer
	cmd.Stdout = writer

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	fmt.Println(buffer.String())
}
