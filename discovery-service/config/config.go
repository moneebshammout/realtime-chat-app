package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port        string
	GatewayHost string
	App         string
	ZooHosts    []string
	SignatureKey string
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
		zooHosts := getEnvVar("ZOO_HOSTS")

		Env = &AppConfig{
			Port:        getEnvVar("PORT"),
			GatewayHost: getEnvVar("GATEWAY_HOST"),
			App:         os.Getenv("App"),
			ZooHosts:    strings.Split(zooHosts, ","),
			SignatureKey: getEnvVar("SIGNATURE_KEY"),
		}
	})
}
