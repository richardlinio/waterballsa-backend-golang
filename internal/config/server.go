package config

import (
	"os"
	"strconv"
	"time"
)

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() ServerConfig {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0" // default to all interfaces
	}

	portStr := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 8080 // default port
	}

	readTimeoutStr := os.Getenv("SERVER_READ_TIMEOUT")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil || readTimeoutStr == "" {
		readTimeout = 10 * time.Second // default
	}

	writeTimeoutStr := os.Getenv("SERVER_WRITE_TIMEOUT")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil || writeTimeoutStr == "" {
		writeTimeout = 10 * time.Second // default
	}

	shutdownTimeoutStr := os.Getenv("SERVER_SHUTDOWN_TIMEOUT")
	shutdownTimeout, err := time.ParseDuration(shutdownTimeoutStr)
	if err != nil || shutdownTimeoutStr == "" {
		shutdownTimeout = 5 * time.Second // default
	}

	requestTimeoutStr := os.Getenv("SERVER_REQUEST_TIMEOUT")
	requestTimeout, err := time.ParseDuration(requestTimeoutStr)
	if err != nil || requestTimeoutStr == "" {
		requestTimeout = 10 * time.Second // default
	}

	return ServerConfig{
		Host:            host,
		Port:            port,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
		RequestTimeout:  requestTimeout,
	}
}
