package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// loadCORSConfig loads CORS configuration from environment variables
func loadCORSConfig() (CORSConfig, error) {
	var config CORSConfig

	// Load allowed origins (comma-separated list)
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsEnv == "" {
		return config, fmt.Errorf("CORS_ALLOWED_ORIGINS is required")
	}
	origins := strings.Split(originsEnv, ",")
	// Trim whitespace and validate each origin
	config.AllowedOrigins = make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			return config, fmt.Errorf("empty origin in CORS_ALLOWED_ORIGINS")
		}
		config.AllowedOrigins = append(config.AllowedOrigins, origin)
	}

	// Load allow credentials (default: true)
	allowCredentialsEnv := os.Getenv("CORS_ALLOW_CREDENTIALS")
	if allowCredentialsEnv == "" {
		config.AllowCredentials = true
	} else {
		var err error
		config.AllowCredentials, err = strconv.ParseBool(allowCredentialsEnv)
		if err != nil {
			return config, fmt.Errorf("invalid CORS_ALLOW_CREDENTIALS value: %w", err)
		}
	}

	// Load max age (default: 3600 seconds)
	maxAgeEnv := os.Getenv("CORS_MAX_AGE_SECONDS")
	if maxAgeEnv == "" {
		config.MaxAge = 3600 * time.Second
	} else {
		maxAgeSeconds, err := strconv.Atoi(maxAgeEnv)
		if err != nil {
			return config, fmt.Errorf("invalid CORS_MAX_AGE_SECONDS value: %w", err)
		}
		config.MaxAge = time.Duration(maxAgeSeconds) * time.Second
	}

	return config, nil
}
