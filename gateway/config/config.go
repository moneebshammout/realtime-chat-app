package config

import (
	"fmt"
	"os"
	"sync"

	"gateway/pkg/utils"
)

var (
	Env           *AppConfig
	Gateway       *GatewayConfig
	appConfigOnce sync.Once
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port            string
	JWTAccessSecret string
}

type ServiceConfig struct {
	Title   string   `json:"title"`
	Paths   []string `json:"paths"`
	Backend string   `json:"backend"`
}

// Config struct to hold the entire Gateway configuration
type GatewayConfig struct {
	Services []ServiceConfig `json:"services"`
	Public   []string        `json:"public"`
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
			Port:            getEnvVar("PORT"),
			JWTAccessSecret: getEnvVar("JWT_ACCESS_SECRET"),
		}

		config, err := utils.ParseJsonFile("gateway_config.json", &GatewayConfig{})
		if err != nil {
			panic(fmt.Sprintf("Error loading gateway_config.json: %s", err.Error()))
		}

		Gateway = config.(*GatewayConfig)
	})
}
