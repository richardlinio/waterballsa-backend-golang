package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// loadDatabaseConfig loads database configuration from environment variables and config file
func loadDatabaseConfig() (*DatabaseConfig, error) {
	// Read from environment variables first, fallback to config.yaml
	host := getEnv("DB_HOST", viper.GetString("database.host"))
	portStr := getEnv("DB_PORT", strconv.Itoa(viper.GetInt("database.port")))
	user := getEnv("DB_USER", viper.GetString("database.user"))
	password := getEnv("DB_PASSWORD", viper.GetString("database.password"))
	name := getEnv("DB_NAME", viper.GetString("database.name"))
	sslmode := getEnv("DB_SSLMODE", viper.GetString("database.sslmode"))
	if sslmode == "" {
		sslmode = "disable" // default to disable for development
	}

	maxConnsStr := getEnv("DB_MAX_CONNS", strconv.Itoa(viper.GetInt("database.max_conns")))
	minConnsStr := getEnv("DB_MIN_CONNS", strconv.Itoa(viper.GetInt("database.min_conns")))

	// Validate required fields
	if host == "" {
		return nil, fmt.Errorf("database host is required (set DB_HOST env var or database.host in config)")
	}
	if user == "" {
		return nil, fmt.Errorf("database user is required (set DB_USER env var or database.user in config)")
	}
	if name == "" {
		return nil, fmt.Errorf("database name is required (set DB_NAME env var or database.name in config)")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432 // default PostgreSQL port
	}

	// Parse pool settings with defaults
	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil || maxConns == 0 {
		maxConns = 10 // default
	}

	minConns, err := strconv.Atoi(minConnsStr)
	if err != nil || minConns == 0 {
		minConns = 2 // default
	}

	config := &DatabaseConfig{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		Name:            name,
		SSLMode:         sslmode,
		MaxConns:        int32(maxConns),
		MinConns:        int32(minConns),
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	return config, nil
}
