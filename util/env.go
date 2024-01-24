package util

import (
	"fmt"
	"go-graphql-api/util/logger"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	_load_env_mu sync.Mutex
	_loaded      = false
)

func _load_env() {
	_load_env_mu.Lock()
	defer _load_env_mu.Unlock()
	if _loaded {
		return
	}

	err := godotenv.Load()
	if err != nil {
		logger.Err("Error loading dotenv: #%v", err)
	} else {
		_loaded = true
	}
}

// Find the value of the environment variable defined with the `envkey`.
// If the variable is not defined, use the provided `default_value`
func EnvOrDefault(envkey string, default_value string) string {
	_load_env()
	value := os.Getenv(envkey)
	if value == "" {
		logger.Err("Env variable \"%s\" not defined, using default value: \"%s\"", envkey, default_value)
		value = default_value
	}
	return value
}

func ServerPort() string {
	return EnvOrDefault("SERVER_PORT", "8080")
}

func ServerUri() string {
	host := EnvOrDefault("SERVER_HOST", "http://localhost")
	return fmt.Sprintf("%s:%s", host, ServerPort())
}
