package app

import (
	"os"
	"sync"

	"relay-service/pkg/utils"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
)
var logger = utils.GetLogger()

// AppConfig holds the application configuration.
type AppConfig struct {
	Port                string
	Host                string
	GatewayHost         string
	App                 string
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Panicf("AppConfig: %s environment variable not set", key)
	}
	return value
}

// init initializes the AppConfig singleton.
func init() {
	appConfigOnce.Do(func() {
		Env = &AppConfig{
			Port:                getEnvVar("PORT"),
			Host:                getEnvVar("HOST"),
			GatewayHost:         getEnvVar("GATEWAY_HOST"),
			App:                 os.Getenv("App"),
		}
	})
}
