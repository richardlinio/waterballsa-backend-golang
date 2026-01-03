package config

import (
	"os"
	"strings"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string
	Format string
}

// loadLoggerConfig loads logger configuration from environment variables
func loadLoggerConfig() LoggerConfig {
	// Get log level from environment (default: info)
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "" {
		level = "info"
	}

	// Get log format from environment (default: text)
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "" {
		format = "text"
	}

	return LoggerConfig{
		Level:  level,
		Format: format,
	}
}
