package logger

import (
	"log/slog"
	"os"

	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

// NewLogger creates and configures a new slog.Logger instance
// based on the provided LoggerConfig
func NewLogger(cfg config.LoggerConfig) *slog.Logger {
	level := parseLogLevel(cfg.Level)

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch cfg.Format {
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
