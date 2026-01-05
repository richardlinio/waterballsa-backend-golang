package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{
			name:     "debug lowercase",
			input:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "debug uppercase",
			input:    "DEBUG",
			expected: slog.LevelDebug,
		},
		{
			name:     "debug mixed case",
			input:    "DeBuG",
			expected: slog.LevelDebug,
		},
		{
			name:     "info lowercase",
			input:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "info uppercase",
			input:    "INFO",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn lowercase",
			input:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "warning full word",
			input:    "warning",
			expected: slog.LevelWarn,
		},
		{
			name:     "warn uppercase",
			input:    "WARN",
			expected: slog.LevelWarn,
		},
		{
			name:     "error lowercase",
			input:    "error",
			expected: slog.LevelError,
		},
		{
			name:     "error uppercase",
			input:    "ERROR",
			expected: slog.LevelError,
		},
		{
			name:     "invalid level defaults to info",
			input:    "invalid",
			expected: slog.LevelInfo,
		},
		{
			name:     "empty string defaults to info",
			input:    "",
			expected: slog.LevelInfo,
		},
		{
			name:     "random string defaults to info",
			input:    "foobar",
			expected: slog.LevelInfo,
		},
		{
			name:     "numeric string defaults to info",
			input:    "123",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

//nolint:gocyclo // Table-driven test with multiple scenarios
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config config.LoggerConfig
		verify func(*testing.T, *slog.Logger)
	}{
		{
			name: "json format with debug level",
			config: config.LoggerConfig{
				Level:  "debug",
				Format: "json",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.Enabled(context.TODO(), slog.LevelDebug) {
					t.Error("logger should be enabled for debug level")
				}
			},
		},
		{
			name: "text format with info level",
			config: config.LoggerConfig{
				Level:  "info",
				Format: "text",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.Enabled(context.TODO(), slog.LevelInfo) {
					t.Error("logger should be enabled for info level")
				}
				if logger.Enabled(context.TODO(), slog.LevelDebug) {
					t.Error("logger should not be enabled for debug level")
				}
			},
		},
		{
			name: "text format with warn level",
			config: config.LoggerConfig{
				Level:  "warn",
				Format: "text",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.Enabled(context.TODO(), slog.LevelWarn) {
					t.Error("logger should be enabled for warn level")
				}
				if logger.Enabled(context.TODO(), slog.LevelInfo) {
					t.Error("logger should not be enabled for info level")
				}
			},
		},
		{
			name: "text format with error level",
			config: config.LoggerConfig{
				Level:  "error",
				Format: "text",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.Enabled(context.TODO(), slog.LevelError) {
					t.Error("logger should be enabled for error level")
				}
				if logger.Enabled(context.TODO(), slog.LevelWarn) {
					t.Error("logger should not be enabled for warn level")
				}
			},
		},
		{
			name: "invalid format defaults to text",
			config: config.LoggerConfig{
				Level:  "info",
				Format: "invalid",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
			},
		},
		{
			name: "invalid level defaults to info",
			config: config.LoggerConfig{
				Level:  "invalid",
				Format: "text",
			},
			verify: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.Enabled(context.TODO(), slog.LevelInfo) {
					t.Error("logger should be enabled for info level")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)
			if tt.verify != nil {
				tt.verify(t, logger)
			}
		})
	}
}
