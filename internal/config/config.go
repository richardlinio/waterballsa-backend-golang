package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Try to read config file, but don't fail if it doesn't exist
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Ignore error if config file doesn't exist
	_ = viper.ReadInConfig()

	// Load database config
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load server config
	serverConfig := loadServerConfig()

	return &Config{
		Server:   serverConfig,
		Database: *dbConfig,
	}, nil
}
