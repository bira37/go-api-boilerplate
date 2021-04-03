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

## Running the API

To run the API, first you need to start the database contained in `docker-compose.yml` file, and then run the migrations. This can be achieved running:

```
docker-compose up -d # Start the CockroachDB container
go run api/cmd/migrate/main.go up # Apply all migrations
go run api/cmd/migrate/main.go up # Start th REST Server
```


## Testing

This project contains both unit and integration tests. For this reason, you also need to start the database and run all migrations before running them. Tests were made using stdlib testing package. To see a detailed coverage visualization, you need to install GoCov CLI. You can execute tests and coverage tool using the following commands:

```
go test ./... # Execute tests
go test -coverprofile coverage.out ./... # Execute tests and generate coverage output also for external packages
go test -coverpkg ./... -coverprofile coverage.out ./... # Execute tests and generate coverage output also for external packages
go tool cover -html=coverage.out -o coverage.html # Generate coverage.html file for better coverage visualization
gocov convert coverage.out | gocov report # Visualize consolidated coverage
```

## TODO

- gRPC sample with tests
