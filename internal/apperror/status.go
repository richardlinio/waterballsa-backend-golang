package apperror

import "net/http"

var httpStatusMap = map[string]int{
	// Validation errors (400)
	CodeValidationFailed: http.StatusBadRequest,
	CodePasswordTooLong:  http.StatusBadRequest,

	// Auth errors (401)
	CodeAuthFailed:   http.StatusUnauthorized, // Login failures, wrong credentials, rate limited
	CodeUnauthorized: http.StatusUnauthorized, // Token invalid/missing

	// Conflict errors (409)
	CodeUsernameExists: http.StatusConflict,

	// Server errors (500)
	CodeInternalError: http.StatusInternalServerError,
	CodeDatabaseError: http.StatusInternalServerError,

	// Service unavailable errors (503)
	CodeServiceUnavailable: http.StatusServiceUnavailable,
}

func GetHTTPStatus(code string) int {
	if status, ok := httpStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError // default to Internal Server Error
}
