package config

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates and configures a new slog.Logger instance
// based on environment variables LOG_LEVEL and LOG_FORMAT
func NewLogger() *slog.Logger {
	// Get log level from environment (default: info)
	levelStr := strings.ToLower(getEnv("LOG_LEVEL", "info"))
	level := parseLogLevel(levelStr)

	// Get log format from environment (default: json)
	formatStr := strings.ToLower(getEnv("LOG_FORMAT", "json"))

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch formatStr {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "json":
		fallthrough
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// parseLogLevel converts string to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
