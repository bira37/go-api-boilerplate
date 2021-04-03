# Go Api Template

This API model was made using:

- [Goose (pressly/goose)](https://github.com/pressly/goose) for migration handling
- [Gin (gin-gonic)](https://github.com/gin-gonic/gin) to construct API handlers and middlewares
- [Sqlx (jmoiron/sqlx)](https://github.com/jmoiron/sqlx) to make the database access layer
- [Testify (stretchr/testify)](https://github.com/stretchr/testify) for mocking dependencies on unit testing
- [GoCov (axw/gocov)](https://github.com/axw/gocov) to generate the consolidated test coverage

## Migrations

Migrations use goose CLI tool. New migrations can be generated executing api/cmd/migrate/main.go

```
go run api/cmd/migrate/main.go up # Upgrade database to the latest migration
go run api/cmd/migrate/main.go create <migration_name> # Create a new empty migration file on api/migrations directory
```

## Testing

Tests were made using stdlib testing package. They can be executed using the following commands:

```
go test ./... # Execute tests
go test -coverpkg ./... -coverprofile coverage.out ./... # Execute tests and generate coverage output also for external packages
go tool cover -html=coverage.out -o coverage.html # Generate coverage.html file for better coverage visualization
gocov convert coverage.out | gocov report # Visualize consolidated coverage
```

## TODO

- Integration tests
