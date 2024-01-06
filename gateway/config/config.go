package config

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

var (
	Env           *AppConfig
	appConfigOnce sync.Once
)

// AppConfig holds the application configuration.
type AppConfig struct {
	Port            string
	JWTAccessSecret string
	UserServiceURL  *url.URL
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable not set", key))
	}
	return value
}

func parseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return parsedURL
}

// init initializes the AppConfig singleton.
func init() {
	appConfigOnce.Do(func() {
		Env = &AppConfig{
			Port:            getEnvVar("PORT"),
			JWTAccessSecret: getEnvVar("JWT_ACCESS_SECRET"),
			UserServiceURL:  parseURL(getEnvVar("USER_SERVICE_URL")),
		}
	})
}
