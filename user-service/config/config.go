package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port                    string
	JWTAccessSecret         string
	JWTRefreshSecret        string
	JWTAccessExpirayMinutes int64
	JWTRefreshExpirayHours  int64
	PostgresAddr            string
	GatewayHost             string
	App                     string
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
		// ParseInt returns two values: the parsed integer and an error
		accessExpiryMinutes, err := strconv.ParseInt(getEnvVar("JWT_ACCESS_EXPIARY_MINUTES"), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing JWT_ACCESS_EXPIARY_MINUTES: %s", err))
		}

		refreshExpiryHours, err := strconv.ParseInt(getEnvVar("JWT_REFRESH_EXPIARY_HOURS"), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing JWT_REFRESH_EXPIARY_HOURS: %s", err))
		}

		Env = &AppConfig{
			Port:                    getEnvVar("PORT"),
			JWTAccessSecret:         getEnvVar("JWT_ACCESS_SECRET"),
			JWTRefreshSecret:        getEnvVar("JWT_REFRESH_SECRET"),
			JWTAccessExpirayMinutes: accessExpiryMinutes,
			JWTRefreshExpirayHours:  refreshExpiryHours,
			PostgresAddr:            getEnvVar("POSTGRES_ADDR"),
			GatewayHost:             getEnvVar("GATEWAY_HOST"),
			App:                     os.Getenv("App"),
		}
	})
}
