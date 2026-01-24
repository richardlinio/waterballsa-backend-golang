package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/richardlinio/waterballsa-backend-golang/internal/config"
)

// NewLogger creates and configures a new slog.Logger instance
// based on the provided LoggerConfig
func NewLogger(cfg config.LoggerConfig) *slog.Logger {
	level := parseLogLevel(cfg.Level)

	// Determine output stream based on log level
	// WARN and ERROR go to stderr, INFO and DEBUG go to stdout
	var output *os.File
	if level >= slog.LevelWarn {
		output = os.Stderr
	} else {
		output = os.Stdout
	}

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	case "text":
		handler = slog.NewTextHandler(output, opts)
	default:
		fmt.Fprintf(os.Stderr, "Warning: invalid log format '%s', using default 'text'\n", cfg.Format)
		handler = slog.NewTextHandler(output, opts)
	}

	return slog.New(handler)
}

// parseLogLevel converts string to slog.Level
func parseLogLevel(level string) slog.Level {
	// Convert to lowercase for case-insensitive comparison
	level = strings.ToLower(level)

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
		fmt.Fprintf(os.Stderr, "Warning: invalid log level '%s', using default 'info'\n", level)
		return slog.LevelInfo
	}
}
