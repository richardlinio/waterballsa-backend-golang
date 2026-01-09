package config

import (
	"os"
	"strconv"
	"time"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	WindowDuration    time.Duration
	CleanupInterval   time.Duration
	RecordExpiry      time.Duration
}

// loadRateLimitConfig loads rate limiting configuration from environment variables
func loadRateLimitConfig() RateLimitConfig {
	enabledStr := os.Getenv("RATE_LIMIT_ENABLED")
	enabled := true // default to enabled
	if enabledStr != "" {
		if parsedEnabled, err := strconv.ParseBool(enabledStr); err == nil {
			enabled = parsedEnabled
		}
	}

	requestsPerMinuteStr := os.Getenv("RATE_LIMIT_REQUESTS_PER_MINUTE")
	requestsPerMinute, err := strconv.Atoi(requestsPerMinuteStr)
	if err != nil || requestsPerMinuteStr == "" || requestsPerMinute <= 0 {
		requestsPerMinute = 60 // default to 60 requests per minute
	}

	windowDurationStr := os.Getenv("RATE_LIMIT_WINDOW_DURATION")
	windowDuration, err := time.ParseDuration(windowDurationStr)
	if err != nil || windowDurationStr == "" {
		windowDuration = 1 * time.Minute // default to 1 minute
	}

	cleanupIntervalStr := os.Getenv("RATE_LIMIT_CLEANUP_INTERVAL")
	cleanupInterval, err := time.ParseDuration(cleanupIntervalStr)
	if err != nil || cleanupIntervalStr == "" {
		cleanupInterval = 5 * time.Minute // default to 5 minutes
	}

	recordExpiryStr := os.Getenv("RATE_LIMIT_RECORD_EXPIRY")
	recordExpiry, err := time.ParseDuration(recordExpiryStr)
	if err != nil || recordExpiryStr == "" {
		recordExpiry = 2 * time.Minute // default to 2 minutes
	}

	return RateLimitConfig{
		Enabled:           enabled,
		RequestsPerMinute: requestsPerMinute,
		WindowDuration:    windowDuration,
		CleanupInterval:   cleanupInterval,
		RecordExpiry:      recordExpiry,
	}
}
