package util

import (
	"go-graphql-api/util/logger"
	"os"
)

// Find the value of the environment variable defined with the `envkey`.
// If the variable is not defined, use the provided `default_value`
func EnvOrDefault(envkey string, default_value string) string {
	value := os.Getenv(envkey)
	if value == "" {
		logger.Err("Env variable \"%s\" not defined, using default value: \"%s\"", envkey, default_value)
		value = default_value
	}
	return value
}
