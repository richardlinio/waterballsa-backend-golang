package config

import (
	"os"
	"testing"
)

const (
	formatJSON = "json"
)

//nolint:gocyclo // Table-driven test with multiple scenarios
func TestLoadLoggerConfig(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"LOG_LEVEL":  os.Getenv("LOG_LEVEL"),
		"LOG_FORMAT": os.Getenv("LOG_FORMAT"),
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
		validate func(*testing.T, LoggerConfig)
	}{
		{
			name:    "defaults when no env vars set",
			envVars: map[string]string{},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != defaultLogLevel {
					t.Errorf("Level = %v; want default info", cfg.Level)
				}
				if cfg.Format != defaultLogFormat {
					t.Errorf("Format = %v; want default text", cfg.Format)
				}
			},
		},
		{
			name: "custom level and format",
			envVars: map[string]string{
				"LOG_LEVEL":  "debug",
				"LOG_FORMAT": "json",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "debug" {
					t.Errorf("Level = %v; want debug", cfg.Level)
				}
				if cfg.Format != formatJSON {
					t.Errorf("Format = %v; want json", cfg.Format)
				}
			},
		},
		{
			name: "uppercase values are converted to lowercase",
			envVars: map[string]string{
				"LOG_LEVEL":  "DEBUG",
				"LOG_FORMAT": "JSON",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "debug" {
					t.Errorf("Level = %v; want debug (lowercased)", cfg.Level)
				}
				if cfg.Format != formatJSON {
					t.Errorf("Format = %v; want json (lowercased)", cfg.Format)
				}
			},
		},
		{
			name: "mixed case values are converted to lowercase",
			envVars: map[string]string{
				"LOG_LEVEL":  "WaRn",
				"LOG_FORMAT": "TeXt",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "warn" {
					t.Errorf("Level = %v; want warn (lowercased)", cfg.Level)
				}
				if cfg.Format != "text" {
					t.Errorf("Format = %v; want text (lowercased)", cfg.Format)
				}
			},
		},
		{
			name: "empty level defaults to info",
			envVars: map[string]string{
				"LOG_LEVEL":  "",
				"LOG_FORMAT": "json",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != defaultLogLevel {
					t.Errorf("Level = %v; want default info", cfg.Level)
				}
				if cfg.Format != formatJSON {
					t.Errorf("Format = %v; want json", cfg.Format)
				}
			},
		},
		{
			name: "empty format defaults to text",
			envVars: map[string]string{
				"LOG_LEVEL":  "error",
				"LOG_FORMAT": "",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "error" {
					t.Errorf("Level = %v; want error", cfg.Level)
				}
				if cfg.Format != defaultLogFormat {
					t.Errorf("Format = %v; want default text", cfg.Format)
				}
			},
		},
		{
			name: "level info and format text",
			envVars: map[string]string{
				"LOG_LEVEL":  "info",
				"LOG_FORMAT": "text",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "info" {
					t.Errorf("Level = %v; want info", cfg.Level)
				}
				if cfg.Format != "text" {
					t.Errorf("Format = %v; want text", cfg.Format)
				}
			},
		},
		{
			name: "level warn",
			envVars: map[string]string{
				"LOG_LEVEL": "warn",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "warn" {
					t.Errorf("Level = %v; want warn", cfg.Level)
				}
				if cfg.Format != defaultLogFormat {
					t.Errorf("Format = %v; want default text", cfg.Format)
				}
			},
		},
		{
			name: "level error",
			envVars: map[string]string{
				"LOG_LEVEL": "error",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "error" {
					t.Errorf("Level = %v; want error", cfg.Level)
				}
				if cfg.Format != defaultLogFormat {
					t.Errorf("Format = %v; want default text", cfg.Format)
				}
			},
		},
		{
			name: "invalid values are accepted (validation happens in logger)",
			envVars: map[string]string{
				"LOG_LEVEL":  "invalid",
				"LOG_FORMAT": "invalid",
			},
			validate: func(t *testing.T, cfg LoggerConfig) {
				if cfg.Level != "invalid" {
					t.Errorf("Level = %v; want invalid", cfg.Level)
				}
				if cfg.Format != "invalid" {
					t.Errorf("Format = %v; want invalid", cfg.Format)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all LOG-related environment variables
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
			cfg := loadLoggerConfig()

			// Validate result
			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
