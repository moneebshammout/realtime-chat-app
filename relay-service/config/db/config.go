package db

import (
	"os"
	"sync"

	"relay-service/pkg/utils"
)

var (
	Env          *DBConfig
	dbConfigOnce sync.Once
)
var logger = utils.GetLogger()

// DBConfig holds the application configuration.
type DBConfig struct {
	ClusterUrl string
	KeySpace   string
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		logger.Panicf("DBConfig: %s environment variable not set", key)
	}
	return value
}

// init initializes the QueueConfig singleton.
func init() {
	dbConfigOnce.Do(func() {
		Env = &DBConfig{
			ClusterUrl: getEnvVar("DATABASE_CLUSTER_HOST"),
			KeySpace:   getEnvVar("DATABASE_KEYSPACE", "relay"),
		}
	})
}
