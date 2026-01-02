package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
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

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() (*DatabaseConfig, error) {
	// Read from environment variables
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable" // default to disable for development
	}

	maxConnsStr := os.Getenv("DB_MAX_CONNS")
	minConnsStr := os.Getenv("DB_MIN_CONNS")

	// Validate required fields
	if host == "" {
		return nil, fmt.Errorf("database host is required (set DB_HOST environment variable)")
	}
	if user == "" {
		return nil, fmt.Errorf("database user is required (set DB_USER environment variable)")
	}
	if name == "" {
		return nil, fmt.Errorf("database name is required (set DB_NAME environment variable)")
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
