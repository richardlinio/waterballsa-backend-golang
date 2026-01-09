package apperror

const (
	// Validation errors (400 - Bad Request)
	CodeValidationFailed = "ERR_VALIDATION_FAILED"
	CodePasswordTooLong  = "ERR_PASSWORD_TOO_LONG"

	// Auth errors (401 - Unauthorized)
	// Note: Use CodeAuthFailed for all authentication failures (login failed, wrong credentials, rate limited)
	// to avoid leaking information about which part of authentication failed
	CodeAuthFailed   = "ERR_AUTH_FAILED"  // Generic auth failure (login, wrong username/password, rate limit)
	CodeUnauthorized = "ERR_UNAUTHORIZED" // Token invalid/missing (for protected endpoints)

	// Conflict errors (409 - Conflict)
	CodeUsernameExists = "ERR_USERNAME_EXISTS"

	// Server errors (500 - Internal Server Error)
	CodeInternalError = "ERR_INTERNAL_ERROR"
	CodeDatabaseError = "ERR_DATABASE_ERROR"

	// Service unavailable errors (503 - Service Unavailable)
	CodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"

	// Rate limiting errors (429 - Too Many Requests)
	CodeRateLimitExceeded = "ERR_RATE_LIMIT_EXCEEDED"
)
