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

	// Not found errors (404 - Not Found)
	CodeJourneyNotFound = "ERR_JOURNEY_NOT_FOUND"
	CodeMissionNotFound = "ERR_MISSION_NOT_FOUND"

	// Rate limiting errors (429 - Too Many Requests)
	CodeRateLimitExceeded = "ERR_RATE_LIMIT_EXCEEDED"

	// Server errors (500 - Internal Server Error)
	CodeInternalServerError = "ERR_INTERNAL_SERVER_ERROR"
	CodeDatabaseError       = "ERR_DATABASE_ERROR"
	CodeAuthStateError      = "ERR_AUTH_STATE_ERROR" // Auth state error (internal logic error in auth flow, should only be used for critical bugs)

	// Service unavailable errors (503 - Service Unavailable)
	CodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"
)
