package config

import (
	"os"
	"testing"
	"time"
)

func TestParseInt32(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal int32
		expected   int32
	}{
		{
			name:       "empty string returns default",
			input:      "",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "valid positive integer",
			input:      "25",
			defaultVal: 10,
			expected:   25,
		},
		{
			name:       "zero returns default",
			input:      "0",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "negative number returns default",
			input:      "-5",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "invalid string returns default",
			input:      "invalid",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "max int32 value",
			input:      "2147483647",
			defaultVal: 10,
			expected:   2147483647,
		},
		{
			name:       "overflow int32 returns default",
			input:      "2147483648",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "very large number returns default",
			input:      "9999999999",
			defaultVal: 10,
			expected:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInt32(tt.input, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("parseInt32(%q, %d) = %d; want %d", tt.input, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal time.Duration
		expected   time.Duration
	}{
		{
			name:       "empty string returns default",
			input:      "",
			defaultVal: 5 * time.Second,
			expected:   5 * time.Second,
		},
		{
			name:       "valid duration in seconds",
			input:      "10s",
			defaultVal: 5 * time.Second,
			expected:   10 * time.Second,
		},
		{
			name:       "valid duration in minutes",
			input:      "2m",
			defaultVal: 5 * time.Second,
			expected:   2 * time.Minute,
		},
		{
			name:       "valid duration in hours",
			input:      "1h",
			defaultVal: 5 * time.Second,
			expected:   1 * time.Hour,
		},
		{
			name:       "zero duration returns default",
			input:      "0s",
			defaultVal: 5 * time.Second,
			expected:   5 * time.Second,
		},
		{
			name:       "negative duration returns default",
			input:      "-10s",
			defaultVal: 5 * time.Second,
			expected:   5 * time.Second,
		},
		{
			name:       "invalid format returns default",
			input:      "invalid",
			defaultVal: 5 * time.Second,
			expected:   5 * time.Second,
		},
		{
			name:       "complex duration",
			input:      "1h30m45s",
			defaultVal: 5 * time.Second,
			expected:   1*time.Hour + 30*time.Minute + 45*time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDuration(tt.input, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("parseDuration(%q, %v) = %v; want %v", tt.input, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

//nolint:gocyclo // Table-driven test with multiple scenarios
func TestLoadDatabaseConfig(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"DB_HOST":               os.Getenv("DB_HOST"),
		"DB_PORT":               os.Getenv("DB_PORT"),
		"DB_USER":               os.Getenv("DB_USER"),
		"DB_PASSWORD":           os.Getenv("DB_PASSWORD"),
		"DB_NAME":               os.Getenv("DB_NAME"),
		"DB_SSLMODE":            os.Getenv("DB_SSLMODE"),
		"DB_MAX_CONNS":          os.Getenv("DB_MAX_CONNS"),
		"DB_MIN_CONNS":          os.Getenv("DB_MIN_CONNS"),
		"DB_MAX_CONN_LIFETIME":  os.Getenv("DB_MAX_CONN_LIFETIME"),
		"DB_MAX_CONN_IDLE_TIME": os.Getenv("DB_MAX_CONN_IDLE_TIME"),
		"DB_CONNECT_TIMEOUT":    os.Getenv("DB_CONNECT_TIMEOUT"),
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
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *DatabaseConfig)
	}{
		{
			name: "valid configuration with all fields",
			envVars: map[string]string{
				"DB_HOST":               "localhost",
				"DB_PORT":               "5432",
				"DB_USER":               "testuser",
				"DB_PASSWORD":           "testpass",
				"DB_NAME":               "testdb",
				"DB_SSLMODE":            "disable",
				"DB_MAX_CONNS":          "20",
				"DB_MIN_CONNS":          "5",
				"DB_MAX_CONN_LIFETIME":  "2h",
				"DB_MAX_CONN_IDLE_TIME": "1h",
				"DB_CONNECT_TIMEOUT":    "30s",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *DatabaseConfig) {
				if cfg.Host != "localhost" {
					t.Errorf("Host = %v; want localhost", cfg.Host)
				}
				if cfg.Port != 5432 {
					t.Errorf("Port = %v; want 5432", cfg.Port)
				}
				if cfg.User != "testuser" {
					t.Errorf("User = %v; want testuser", cfg.User)
				}
				if cfg.Password != "testpass" {
					t.Errorf("Password = %v; want testpass", cfg.Password)
				}
				if cfg.Name != "testdb" {
					t.Errorf("Name = %v; want testdb", cfg.Name)
				}
				if cfg.SSLMode != defaultSSLMode {
					t.Errorf("SSLMode = %v; want disable", cfg.SSLMode)
				}
				if cfg.MaxConns != 20 {
					t.Errorf("MaxConns = %v; want 20", cfg.MaxConns)
				}
				if cfg.MinConns != 5 {
					t.Errorf("MinConns = %v; want 5", cfg.MinConns)
				}
				if cfg.MaxConnLifetime != 2*time.Hour {
					t.Errorf("MaxConnLifetime = %v; want 2h", cfg.MaxConnLifetime)
				}
				if cfg.MaxConnIdleTime != 1*time.Hour {
					t.Errorf("MaxConnIdleTime = %v; want 1h", cfg.MaxConnIdleTime)
				}
				if cfg.ConnectTimeout != 30*time.Second {
					t.Errorf("ConnectTimeout = %v; want 30s", cfg.ConnectTimeout)
				}
			},
		},
		{
			name: "valid configuration with minimal fields and defaults",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *DatabaseConfig) {
				if cfg.Port != 5432 {
					t.Errorf("Port = %v; want default 5432", cfg.Port)
				}
				if cfg.SSLMode != defaultSSLMode {
					t.Errorf("SSLMode = %v; want default disable", cfg.SSLMode)
				}
				if cfg.MaxConns != 10 {
					t.Errorf("MaxConns = %v; want default 10", cfg.MaxConns)
				}
				if cfg.MinConns != 2 {
					t.Errorf("MinConns = %v; want default 2", cfg.MinConns)
				}
				if cfg.MaxConnLifetime != 1*time.Hour {
					t.Errorf("MaxConnLifetime = %v; want default 1h", cfg.MaxConnLifetime)
				}
				if cfg.MaxConnIdleTime != 30*time.Minute {
					t.Errorf("MaxConnIdleTime = %v; want default 30m", cfg.MaxConnIdleTime)
				}
				if cfg.ConnectTimeout != 10*time.Second {
					t.Errorf("ConnectTimeout = %v; want default 10s", cfg.ConnectTimeout)
				}
			},
		},
		{
			name: "missing DB_HOST returns error",
			envVars: map[string]string{
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "database host is required",
		},
		{
			name: "missing DB_USER returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "database user is required",
		},
		{
			name: "missing DB_PASSWORD returns error",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "testuser",
				"DB_NAME": "testdb",
			},
			wantErr: true,
			errMsg:  "database password is required",
		},
		{
			name: "missing DB_NAME returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
			},
			wantErr: true,
			errMsg:  "database name is required",
		},
		{
			name: "invalid DB_PORT returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "invalid",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "invalid DB_PORT",
		},
		{
			name: "DB_PORT out of range (too low) returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "0",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "invalid DB_PORT: 0 (must be 1-65535)",
		},
		{
			name: "DB_PORT out of range (too high) returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "65536",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "invalid DB_PORT: 65536 (must be 1-65535)",
		},
		{
			name: "invalid DB_SSLMODE returns error",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_SSLMODE":  "invalid",
			},
			wantErr: true,
			errMsg:  "invalid DB_SSLMODE",
		},
		{
			name: "valid DB_SSLMODE require",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_SSLMODE":  "require",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *DatabaseConfig) {
				if cfg.SSLMode != "require" {
					t.Errorf("SSLMode = %v; want require", cfg.SSLMode)
				}
			},
		},
		{
			name: "valid DB_SSLMODE verify-full",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_SSLMODE":  "verify-full",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *DatabaseConfig) {
				if cfg.SSLMode != "verify-full" {
					t.Errorf("SSLMode = %v; want verify-full", cfg.SSLMode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all DB-related environment variables
			for key := range originalEnv {
				if err := os.Unsetenv(key); err != nil {
					t.Fatalf("failed to unset env var %s: %v", key, err)
				}
			}

			// Set test environment variables
			for key, val := range tt.envVars {
				if err := os.Setenv(key, val); err != nil {
					t.Fatalf("failed to set env var %s: %v", key, err)
				}
			}

			// Run test
			cfg, err := loadDatabaseConfig()

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("loadDatabaseConfig() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("loadDatabaseConfig() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("loadDatabaseConfig() unexpected error = %v", err)
				return
			}

			// Validate result
			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
