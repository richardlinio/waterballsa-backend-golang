package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadCORSConfig_ValidSingleOrigin(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS":   "https://example.com",
		"CORS_ALLOW_CREDENTIALS": "true",
		"CORS_MAX_AGE_SECONDS":   "7200",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.AllowedOrigins) != 1 {
		t.Errorf("expected 1 origin, got %d", len(config.AllowedOrigins))
	}
	if config.AllowedOrigins[0] != "https://example.com" {
		t.Errorf("expected origin https://example.com, got %s", config.AllowedOrigins[0])
	}
	if !config.AllowCredentials {
		t.Error("expected AllowCredentials to be true")
	}
	if config.MaxAge != 7200*time.Second {
		t.Errorf("expected MaxAge 7200s, got %v", config.MaxAge)
	}
}

func TestLoadCORSConfig_ValidMultipleOrigins(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS":   "https://example.com,https://api.example.com,http://localhost:3000",
		"CORS_ALLOW_CREDENTIALS": "false",
		"CORS_MAX_AGE_SECONDS":   "1800",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.AllowedOrigins) != 3 {
		t.Errorf("expected 3 origins, got %d", len(config.AllowedOrigins))
	}
	expectedOrigins := []string{
		"https://example.com",
		"https://api.example.com",
		"http://localhost:3000",
	}
	for i, expected := range expectedOrigins {
		if config.AllowedOrigins[i] != expected {
			t.Errorf("origin[%d]: expected %s, got %s", i, expected, config.AllowedOrigins[i])
		}
	}
	if config.AllowCredentials {
		t.Error("expected AllowCredentials to be false")
	}
	if config.MaxAge != 1800*time.Second {
		t.Errorf("expected MaxAge 1800s, got %v", config.MaxAge)
	}
}

func TestLoadCORSConfig_OriginsWithWhitespaceTrimmed(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "  https://example.com  ,  https://api.example.com  ",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.AllowedOrigins) != 2 {
		t.Errorf("expected 2 origins, got %d", len(config.AllowedOrigins))
	}
	if config.AllowedOrigins[0] != "https://example.com" {
		t.Errorf("expected origin https://example.com, got %s", config.AllowedOrigins[0])
	}
	if config.AllowedOrigins[1] != "https://api.example.com" {
		t.Errorf("expected origin https://api.example.com, got %s", config.AllowedOrigins[1])
	}
}

func TestLoadCORSConfig_WildcardOrigin(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://*.example.com",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.AllowedOrigins) != 1 {
		t.Errorf("expected 1 origin, got %d", len(config.AllowedOrigins))
	}
	if config.AllowedOrigins[0] != "https://*.example.com" {
		t.Errorf("expected origin https://*.example.com, got %s", config.AllowedOrigins[0])
	}
}

func TestLoadCORSConfig_DefaultValues(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !config.AllowCredentials {
		t.Error("expected default AllowCredentials to be true")
	}
	if config.MaxAge != 3600*time.Second {
		t.Errorf("expected default MaxAge 3600s, got %v", config.MaxAge)
	}
}

func TestLoadCORSConfig_NegativeMaxAge(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com",
		"CORS_MAX_AGE_SECONDS": "-100",
	})

	config, err := loadCORSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.MaxAge != -100*time.Second {
		t.Errorf("expected MaxAge -100s, got %v", config.MaxAge)
	}
}

func TestLoadCORSConfig_MissingAllowedOrigins(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "CORS_ALLOWED_ORIGINS is required") {
		t.Errorf("expected error containing 'CORS_ALLOWED_ORIGINS is required', got %q", err.Error())
	}
}

func TestLoadCORSConfig_EmptyAllowedOrigins(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "CORS_ALLOWED_ORIGINS is required") {
		t.Errorf("expected error containing 'CORS_ALLOWED_ORIGINS is required', got %q", err.Error())
	}
}

func TestLoadCORSConfig_EmptyOriginAfterSplit(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com,,https://api.example.com",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "empty origin in CORS_ALLOWED_ORIGINS") {
		t.Errorf("expected error containing 'empty origin in CORS_ALLOWED_ORIGINS', got %q", err.Error())
	}
}

func TestLoadCORSConfig_EmptyOriginWithWhitespace(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com,   ,https://api.example.com",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "empty origin in CORS_ALLOWED_ORIGINS") {
		t.Errorf("expected error containing 'empty origin in CORS_ALLOWED_ORIGINS', got %q", err.Error())
	}
}

func TestLoadCORSConfig_TrailingComma(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com,",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "empty origin in CORS_ALLOWED_ORIGINS") {
		t.Errorf("expected error containing 'empty origin in CORS_ALLOWED_ORIGINS', got %q", err.Error())
	}
}

func TestLoadCORSConfig_LeadingComma(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": ",https://example.com",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "empty origin in CORS_ALLOWED_ORIGINS") {
		t.Errorf("expected error containing 'empty origin in CORS_ALLOWED_ORIGINS', got %q", err.Error())
	}
}

func TestLoadCORSConfig_InvalidAllowCredentials(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS":   "https://example.com",
		"CORS_ALLOW_CREDENTIALS": "invalid",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid CORS_ALLOW_CREDENTIALS value") {
		t.Errorf("expected error containing 'invalid CORS_ALLOW_CREDENTIALS value', got %q", err.Error())
	}
}

func TestLoadCORSConfig_InvalidMaxAgeSeconds(t *testing.T) {
	savedEnv := saveAndClearEnv(t)
	defer restoreEnv(t, savedEnv)

	setEnvVars(t, map[string]string{
		"CORS_ALLOWED_ORIGINS": "https://example.com",
		"CORS_MAX_AGE_SECONDS": "not-a-number",
	})

	_, err := loadCORSConfig()
	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid CORS_MAX_AGE_SECONDS value") {
		t.Errorf("expected error containing 'invalid CORS_MAX_AGE_SECONDS value', got %q", err.Error())
	}
}

// saveAndClearEnv saves current CORS environment variables and clears them
func saveAndClearEnv(t *testing.T) map[string]string {
	t.Helper()

	corsEnvVars := []string{
		"CORS_ALLOWED_ORIGINS",
		"CORS_ALLOW_CREDENTIALS",
		"CORS_MAX_AGE_SECONDS",
	}

	savedEnv := make(map[string]string)
	for _, key := range corsEnvVars {
		if val, exists := os.LookupEnv(key); exists {
			savedEnv[key] = val
		}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset env var %s: %v", key, err)
		}
	}

	return savedEnv
}

// restoreEnv restores environment variables from saved state
func restoreEnv(t *testing.T, savedEnv map[string]string) {
	t.Helper()

	// First unset all CORS variables
	corsEnvVars := []string{
		"CORS_ALLOWED_ORIGINS",
		"CORS_ALLOW_CREDENTIALS",
		"CORS_MAX_AGE_SECONDS",
	}
	for _, key := range corsEnvVars {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset env var %s: %v", key, err)
		}
	}

	// Then restore saved values
	for key, val := range savedEnv {
		if err := os.Setenv(key, val); err != nil {
			t.Fatalf("failed to restore env var %s: %v", key, err)
		}
	}
}

// setEnvVars sets environment variables for testing
func setEnvVars(t *testing.T, envVars map[string]string) {
	t.Helper()

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("failed to set env var %s: %v", key, err)
		}
	}
}
