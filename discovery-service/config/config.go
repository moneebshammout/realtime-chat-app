package config

import (
	"os"
	"strings"
	"sync"

	"discovery-service/pkg/utils"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
	logger        = utils.GetLogger()
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port         string
	GatewayPort  string
	GatewayHost  string
	App          string
	ZooHosts     []string
	SignatureKey string
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Panicf("%s environment variable not set", key)
	}
	return value
}

// init initializes the AppConfig singleton.
func init() {
	appConfigOnce.Do(func() {
		zooHosts := getEnvVar("ZOO_HOSTS")

		Env = &AppConfig{
			Port:         getEnvVar("PORT"),
			GatewayPort:  getEnvVar("GATEWAY_PORT"),
			GatewayHost:  getEnvVar("GATEWAY_HOST"),
			App:          os.Getenv("App"),
			ZooHosts:     strings.Split(zooHosts, ","),
			SignatureKey: getEnvVar("SIGNATURE_KEY"),
		}
	})
}
