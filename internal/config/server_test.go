package config

import (
	"os"
	"testing"
	"time"
)

//nolint:gocyclo // Table-driven test with multiple scenarios
func TestLoadServerConfig(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"SERVER_HOST":             os.Getenv("SERVER_HOST"),
		"SERVER_PORT":             os.Getenv("SERVER_PORT"),
		"SERVER_READ_TIMEOUT":     os.Getenv("SERVER_READ_TIMEOUT"),
		"SERVER_WRITE_TIMEOUT":    os.Getenv("SERVER_WRITE_TIMEOUT"),
		"SERVER_SHUTDOWN_TIMEOUT": os.Getenv("SERVER_SHUTDOWN_TIMEOUT"),
		"SERVER_REQUEST_TIMEOUT":  os.Getenv("SERVER_REQUEST_TIMEOUT"),
	}

	// Restore environment after test
	defer func() {
		for key, val := range originalEnv {
			if val == "" {
				if err := os.Unsetenv(key); err != nil {
					t.Fatalf("failed to unset env var %s: %v", key, err)
				}
			} else {
				if err := os.Setenv(key, val); err != nil {
					t.Fatalf("failed to restore env var %s: %v", key, err)
				}
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(*testing.T, ServerConfig)
	}{
		{
			name:    "all defaults when no env vars set",
			envVars: map[string]string{},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.Host != "0.0.0.0" {
					t.Errorf("Host = %v; want default 0.0.0.0", cfg.Host)
				}
				if cfg.Port != 8080 {
					t.Errorf("Port = %v; want default 8080", cfg.Port)
				}
				if cfg.ReadTimeout != 10*time.Second {
					t.Errorf("ReadTimeout = %v; want default 10s", cfg.ReadTimeout)
				}
				if cfg.WriteTimeout != 10*time.Second {
					t.Errorf("WriteTimeout = %v; want default 10s", cfg.WriteTimeout)
				}
				if cfg.ShutdownTimeout != 5*time.Second {
					t.Errorf("ShutdownTimeout = %v; want default 5s", cfg.ShutdownTimeout)
				}
				if cfg.RequestTimeout != 10*time.Second {
					t.Errorf("RequestTimeout = %v; want default 10s", cfg.RequestTimeout)
				}
			},
		},
		{
			name: "custom host and port",
			envVars: map[string]string{
				"SERVER_HOST": "127.0.0.1",
				"SERVER_PORT": "3000",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.Host != "127.0.0.1" {
					t.Errorf("Host = %v; want 127.0.0.1", cfg.Host)
				}
				if cfg.Port != 3000 {
					t.Errorf("Port = %v; want 3000", cfg.Port)
				}
			},
		},
		{
			name: "invalid port defaults to 8080",
			envVars: map[string]string{
				"SERVER_PORT": "invalid",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.Port != 8080 {
					t.Errorf("Port = %v; want default 8080", cfg.Port)
				}
			},
		},
		{
			name: "port zero defaults to 8080",
			envVars: map[string]string{
				"SERVER_PORT": "0",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.Port != 8080 {
					t.Errorf("Port = %v; want default 8080", cfg.Port)
				}
			},
		},
		{
			name: "custom timeouts",
			envVars: map[string]string{
				"SERVER_READ_TIMEOUT":     "30s",
				"SERVER_WRITE_TIMEOUT":    "45s",
				"SERVER_SHUTDOWN_TIMEOUT": "15s",
				"SERVER_REQUEST_TIMEOUT":  "20s",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.ReadTimeout != 30*time.Second {
					t.Errorf("ReadTimeout = %v; want 30s", cfg.ReadTimeout)
				}
				if cfg.WriteTimeout != 45*time.Second {
					t.Errorf("WriteTimeout = %v; want 45s", cfg.WriteTimeout)
				}
				if cfg.ShutdownTimeout != 15*time.Second {
					t.Errorf("ShutdownTimeout = %v; want 15s", cfg.ShutdownTimeout)
				}
				if cfg.RequestTimeout != 20*time.Second {
					t.Errorf("RequestTimeout = %v; want 20s", cfg.RequestTimeout)
				}
			},
		},
		{
			name: "complex duration formats",
			envVars: map[string]string{
				"SERVER_READ_TIMEOUT":  "1m30s",
				"SERVER_WRITE_TIMEOUT": "2h",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.ReadTimeout != 1*time.Minute+30*time.Second {
					t.Errorf("ReadTimeout = %v; want 1m30s", cfg.ReadTimeout)
				}
				if cfg.WriteTimeout != 2*time.Hour {
					t.Errorf("WriteTimeout = %v; want 2h", cfg.WriteTimeout)
				}
			},
		},
		{
			name: "invalid read timeout defaults to 10s",
			envVars: map[string]string{
				"SERVER_READ_TIMEOUT": "invalid",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.ReadTimeout != 10*time.Second {
					t.Errorf("ReadTimeout = %v; want default 10s", cfg.ReadTimeout)
				}
			},
		},
		{
			name: "invalid write timeout defaults to 10s",
			envVars: map[string]string{
				"SERVER_WRITE_TIMEOUT": "invalid",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.WriteTimeout != 10*time.Second {
					t.Errorf("WriteTimeout = %v; want default 10s", cfg.WriteTimeout)
				}
			},
		},
		{
			name: "invalid shutdown timeout defaults to 5s",
			envVars: map[string]string{
				"SERVER_SHUTDOWN_TIMEOUT": "invalid",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.ShutdownTimeout != 5*time.Second {
					t.Errorf("ShutdownTimeout = %v; want default 5s", cfg.ShutdownTimeout)
				}
			},
		},
		{
			name: "invalid request timeout defaults to 10s",
			envVars: map[string]string{
				"SERVER_REQUEST_TIMEOUT": "invalid",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.RequestTimeout != 10*time.Second {
					t.Errorf("RequestTimeout = %v; want default 10s", cfg.RequestTimeout)
				}
			},
		},
		{
			name: "empty timeout strings use defaults",
			envVars: map[string]string{
				"SERVER_READ_TIMEOUT":     "",
				"SERVER_WRITE_TIMEOUT":    "",
				"SERVER_SHUTDOWN_TIMEOUT": "",
				"SERVER_REQUEST_TIMEOUT":  "",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.ReadTimeout != 10*time.Second {
					t.Errorf("ReadTimeout = %v; want default 10s", cfg.ReadTimeout)
				}
				if cfg.WriteTimeout != 10*time.Second {
					t.Errorf("WriteTimeout = %v; want default 10s", cfg.WriteTimeout)
				}
				if cfg.ShutdownTimeout != 5*time.Second {
					t.Errorf("ShutdownTimeout = %v; want default 5s", cfg.ShutdownTimeout)
				}
				if cfg.RequestTimeout != 10*time.Second {
					t.Errorf("RequestTimeout = %v; want default 10s", cfg.RequestTimeout)
				}
			},
		},
		{
			name: "all custom values",
			envVars: map[string]string{
				"SERVER_HOST":             "192.168.1.1",
				"SERVER_PORT":             "9000",
				"SERVER_READ_TIMEOUT":     "25s",
				"SERVER_WRITE_TIMEOUT":    "35s",
				"SERVER_SHUTDOWN_TIMEOUT": "8s",
				"SERVER_REQUEST_TIMEOUT":  "15s",
			},
			validate: func(t *testing.T, cfg ServerConfig) {
				if cfg.Host != "192.168.1.1" {
					t.Errorf("Host = %v; want 192.168.1.1", cfg.Host)
				}
				if cfg.Port != 9000 {
					t.Errorf("Port = %v; want 9000", cfg.Port)
				}
				if cfg.ReadTimeout != 25*time.Second {
					t.Errorf("ReadTimeout = %v; want 25s", cfg.ReadTimeout)
				}
				if cfg.WriteTimeout != 35*time.Second {
					t.Errorf("WriteTimeout = %v; want 35s", cfg.WriteTimeout)
				}
				if cfg.ShutdownTimeout != 8*time.Second {
					t.Errorf("ShutdownTimeout = %v; want 8s", cfg.ShutdownTimeout)
				}
				if cfg.RequestTimeout != 15*time.Second {
					t.Errorf("RequestTimeout = %v; want 15s", cfg.RequestTimeout)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all SERVER-related environment variables
			for key := range originalEnv {
				if err := os.Unsetenv(key); err != nil {
					t.Fatalf("failed to unset env var %s: %v", key, err)
				}
			}

			// Set test environment variables
			for key, val := range tt.envVars {
				if val != "" {
					if err := os.Setenv(key, val); err != nil {
						t.Fatalf("failed to set env var %s: %v", key, err)
					}
				}
			}

			// Run test
			cfg := loadServerConfig()

			// Validate result
			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
