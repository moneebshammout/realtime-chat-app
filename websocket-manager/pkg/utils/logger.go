package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Init initializes the logger with color support and other settings.
func InitLogger() *log.Logger{
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true, // Adds timestamps to each log entry
	})

	log.SetOutput(os.Stdout)    // Log to stdout instead of the default stderr
	log.SetLevel(log.InfoLevel) // Set log level to Info (change to Debug, Warn, Error, etc. as needed)

	return log.StandardLogger()
}

// GetLogger returns the logger instance.
func GetLogger() *log.Logger {
	return log.StandardLogger()
}
