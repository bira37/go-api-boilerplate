package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	SQLDBConnectionString string
	JwtSigningSecret      string
	JwtExpirationMinutes  int
}

var Configuration Config

var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		Configuration = Config{
			SQLDBConnectionString: StringEnvOrDefault("SQLDB_CONNECTION_STRING", "postgres://root@localhost:26257?sslmode=disable"),
			JwtSigningSecret:      StringEnvOrDefault("JWT_SIGNING_SECRET", "da62f84d-e8cc-4725-8725-48f90657ac6c"),
			JwtExpirationMinutes:  86400,
		}
	})
	return Configuration
}

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
