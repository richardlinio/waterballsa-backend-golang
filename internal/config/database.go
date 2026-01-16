package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultSSLMode = "disable"
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
	ConnectTimeout  time.Duration
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() (DatabaseConfig, error) {
	// Read from environment variables
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	// Parse and validate SSL mode
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = defaultSSLMode // default to disable for development
	} else {
		// Validate SSL mode
		validModes := map[string]struct{}{
			defaultSSLMode: {},
			"allow":        {},
			"prefer":       {},
			"require":      {},
			"verify-ca":    {},
			"verify-full":  {},
		}
		if _, ok := validModes[sslmode]; !ok {
			return DatabaseConfig{}, fmt.Errorf("invalid DB_SSLMODE: %s (must be one of: disable, allow, prefer, require, verify-ca, verify-full)", sslmode)
		}
	}

	maxConnsStr := os.Getenv("DB_MAX_CONNS")
	minConnsStr := os.Getenv("DB_MIN_CONNS")
	maxConnLifetimeStr := os.Getenv("DB_MAX_CONN_LIFETIME")
	maxConnIdleTimeStr := os.Getenv("DB_MAX_CONN_IDLE_TIME")
	connectTimeoutStr := os.Getenv("DB_CONNECT_TIMEOUT")

	// Validate required fields
	if host == "" {
		return DatabaseConfig{}, fmt.Errorf("database host is required (set DB_HOST environment variable)")
	}
	if user == "" {
		return DatabaseConfig{}, fmt.Errorf("database user is required (set DB_USER environment variable)")
	}
	if password == "" {
		return DatabaseConfig{}, fmt.Errorf("database password is required (set DB_PASSWORD environment variable)")
	}
	if name == "" {
		return DatabaseConfig{}, fmt.Errorf("database name is required (set DB_NAME environment variable)")
	}

	// Parse port with validation
	port := 5432 // default PostgreSQL port
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return DatabaseConfig{}, fmt.Errorf("invalid DB_PORT: %s (not a number)", portStr)
		}
		if p <= 0 || p > 65535 {
			return DatabaseConfig{}, fmt.Errorf("invalid DB_PORT: %d (must be 1-65535)", p)
		}
		port = p
	}

	// Parse pool settings with defaults
	maxConns := parseInt32(maxConnsStr, 10)
	minConns := parseInt32(minConnsStr, 2)

	// Parse connection lifetime settings with defaults
	maxConnLifetime := parseDuration(maxConnLifetimeStr, 1*time.Hour)
	maxConnIdleTime := parseDuration(maxConnIdleTimeStr, 30*time.Minute)
	connectTimeout := parseDuration(connectTimeoutStr, 10*time.Second)

	return DatabaseConfig{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		Name:            name,
		SSLMode:         sslmode,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
		ConnectTimeout:  connectTimeout,
	}, nil
}

func parseInt32(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return defaultVal
	}
	// Validate int32 range to prevent overflow
	if v > 2147483647 {
		return defaultVal
	}
	return int32(v) // #nosec G109 -- validated above
}

func parseDuration(s string, defaultVal time.Duration) time.Duration {
	if s == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return defaultVal
	}
	return d
}
