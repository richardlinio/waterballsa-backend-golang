package config

import (
	"fmt"
	"log/slog"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   *slog.Logger
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load database config
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load server config
	serverConfig := loadServerConfig()

	// Initialize logger
	logger := NewLogger()

	return &Config{
		Server:   serverConfig,
		Database: *dbConfig,
		Logger:   logger,
	}, nil
}
