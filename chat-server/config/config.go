package config

import (
	"fmt"
	"os"
	"sync"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port string
	App  string
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable not set", key))
	}
	return value
}

// init initializes the AppConfig singleton.
func init() {
	appConfigOnce.Do(func() {
		Env = &AppConfig{
			Port: getEnvVar("PORT"),
			App:  os.Getenv("APP"),
		}
	})
}
