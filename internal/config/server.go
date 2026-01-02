package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// loadServerConfig loads server configuration from environment variables and config file
func loadServerConfig() ServerConfig {
	host := getEnv("SERVER_HOST", viper.GetString("server.host"))
	if host == "" {
		host = "0.0.0.0" // default to all interfaces
	}

	portStr := getEnv("SERVER_PORT", strconv.Itoa(viper.GetInt("server.port")))
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 8080 // default port
	}

	readTimeoutStr := getEnv("SERVER_READ_TIMEOUT", viper.GetString("server.read_timeout"))
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		readTimeout = 10 * time.Second // default
	}

	writeTimeoutStr := getEnv("SERVER_WRITE_TIMEOUT", viper.GetString("server.write_timeout"))
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		writeTimeout = 10 * time.Second // default
	}

	shutdownTimeoutStr := getEnv("SERVER_SHUTDOWN_TIMEOUT", viper.GetString("server.shutdown_timeout"))
	shutdownTimeout, err := time.ParseDuration(shutdownTimeoutStr)
	if err != nil {
		shutdownTimeout = 5 * time.Second // default
	}

	return ServerConfig{
		Host:            host,
		Port:            port,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
