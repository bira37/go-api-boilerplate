package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringEnvOrDefault(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_STRING_ENV", "string")

	if err != nil {
		panic(err)
	}

	value := StringEnvOrDefault("TEST_STRING_ENV", "string2")
	value2 := StringEnvOrDefault("TEST_STRING_ENV_2", "string2")

	assert.Equal(value, "string")
	assert.Equal(value2, "string2")
}

func TestIntEnvOrDefault(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_INT_ENV", "5")

	if err != nil {
		panic(err)
	}

	value := IntEnvOrDefault("TEST_INT_ENV", 10)
	value2 := IntEnvOrDefault("TEST_INT_ENV_2", 10)

	assert.Equal(value, 5)
	assert.Equal(value2, 10)
}

func TestIntEnvOrDefaultPanic(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_INT_ENV_PANIC", "9223372036854775807")

	if err != nil {
		panic(err)
	}

	assert.Panics(func() { IntEnvOrDefault("TEST_INT_ENV_PANIC", 10) })
}

func TestInt64EnvOrDefault(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_INT64_ENV", "9223372036854775807")

	if err != nil {
		panic(err)
	}

	value := Int64EnvOrDefault("TEST_INT64_ENV", 10)
	value2 := Int64EnvOrDefault("TEST_INT64_ENV_2", 10)

	assert.Equal(value, int64(9223372036854775807))
	assert.Equal(value2, int64(10))
}

func TestInt64EnvOrDefaultPanics(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_INT64_ENV_PANIC", "92233720368547758070")

	if err != nil {
		panic(err)
	}

	assert.Panics(func() { Int64EnvOrDefault("TEST_INT64_ENV_PANIC", 10) })
}

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(GetConfig(), GetConfig())
}
