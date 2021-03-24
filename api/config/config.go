package config

import (
	"os"
	"strconv"
)

var (
	SQLDBConnectionString = StringEnvOrDefault("SQLDB_CONNECTION_STRING", "postgres://root@localhost:26257?sslmode=disable")
	JwtSigningSecret      = StringEnvOrDefault("JWT_SIGNING_SECRET", "da62f84d-e8cc-4725-8725-48f90657ac6c")
	JwtExpirationMinute   = 86400
)

func StringEnvOrDefault(env string, def string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return def
}

func IntEnvOrDefault(env string, def int) int {
	val, err := strconv.ParseInt(StringEnvOrDefault(env, strconv.Itoa(def)), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(val)
}

func Int64EnvOrDefault(env string, def int64) int64 {
	val, err := strconv.ParseInt(StringEnvOrDefault(env, strconv.FormatInt(def, 10)), 10, 64)
	if err != nil {
		panic(err)
	}
	return val
}
