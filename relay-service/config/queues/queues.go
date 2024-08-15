package queues

import (
	"os"
	"sync"

	"relay-service/pkg/utils"
)

var (
	Env             *QueueConfig
	queueConfigOnce sync.Once
)
var logger = utils.GetLogger()

// QueueConfig holds the application configuration.
type QueueConfig struct {
	Port              string
	Host              string
	RedisAddr         string
	RelayQueue        string
	MetricsPort       string
	PrometheusAddress string
}

// getEnvVar retrieves an environment variable and returns its value or panics if it's not set.
func getEnvVar(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		logger.Panicf("QueueConfig: %s environment variable not set", key)
	}
	return value
}

// init initializes the QueueConfig singleton.
func init() {
	queueConfigOnce.Do(func() {
		Env = &QueueConfig{
			Port:              getEnvVar("QUEUES_SERVER_PORT"),
			Host:              getEnvVar("QUEUES_SERVER_HOST"),
			RedisAddr:         getEnvVar("REDIS_ADDR"),
			RelayQueue:        getEnvVar("MESSAGE_QUEUE", "relay_queue"),
			MetricsPort:       getEnvVar("METRICS_PORT"),
			PrometheusAddress: getEnvVar("PROMETHEUS_ADDRESS"),
		}
	})
}
