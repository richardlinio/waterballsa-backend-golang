package config

import (
	"fmt"
	"os"
	"time"
)

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret              []byte        // JWT signing secret key
	AccessTokenTimeout  time.Duration // Access token validity duration
	RefreshTokenTimeout time.Duration // Refresh token validity duration
	SecureCookie        bool          // Enable secure cookie (HTTPS only)
	CookieDomain        string        // Cookie domain (optional)
	Realm               string        // Bearer realm for authentication
}

// loadJWTConfig loads JWT configuration from environment variables
func loadJWTConfig() (JWTConfig, error) {
	// Load required JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return JWTConfig{}, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Validate secret length (at least 32 bytes for security)
	if len(secret) < 32 {
		return JWTConfig{}, fmt.Errorf("JWT_SECRET must be at least 32 characters long for security (current: %d)", len(secret))
	}

	// Load optional access token timeout (default: 15 minutes)
	accessTokenTimeout, err := getDurationWithDefault("JWT_ACCESS_TOKEN_TIMEOUT", 15*time.Minute)
	if err != nil {
		return JWTConfig{}, err
	}

	// Load optional refresh token timeout (default: 7 days)
	refreshTokenTimeout, err := getDurationWithDefault("JWT_REFRESH_TOKEN_TIMEOUT", 168*time.Hour)
	if err != nil {
		return JWTConfig{}, err
	}

	// Load optional secure cookie setting (default: false)
	secureCookie := getBoolWithDefault("JWT_SECURE_COOKIE", false)

	// Load optional cookie domain
	cookieDomain := os.Getenv("JWT_COOKIE_DOMAIN")

	// Load optional realm (default: "waterballsa")
	realm := os.Getenv("JWT_REALM")
	if realm == "" {
		realm = "API"
	}

	return JWTConfig{
		Secret:              []byte(secret),
		AccessTokenTimeout:  accessTokenTimeout,
		RefreshTokenTimeout: refreshTokenTimeout,
		SecureCookie:        secureCookie,
		CookieDomain:        cookieDomain,
		Realm:               realm,
	}, nil
}

// getDurationWithDefault loads a duration from environment variable with a default value
func getDurationWithDefault(key string, defaultValue time.Duration) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue, nil
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return duration, nil
}

// getBoolWithDefault loads a boolean from environment variable with a default value
func getBoolWithDefault(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val == "true"
}
