package config

import (
	"fmt"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	CORS     CORSConfig
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

	// Load logger config
	loggerConfig := loadLoggerConfig()

	// Load CORS config
	corsConfig, err := loadCORSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load CORS config: %w", err)
	}

	return &Config{
		Server:   serverConfig,
		Database: *dbConfig,
		Logger:   loggerConfig,
		CORS:     corsConfig,
	}, nil
}
