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

	// Not found errors (404)
	CodeJourneyNotFound: http.StatusNotFound,
	CodeMissionNotFound: http.StatusNotFound,

	// Rate limiting errors (429)
	CodeRateLimitExceeded: http.StatusTooManyRequests,

	// Server errors (500)
	CodeInternalServerError: http.StatusInternalServerError,
	CodeDatabaseError:       http.StatusInternalServerError,
	CodeAuthStateError:      http.StatusInternalServerError, // Auth state error (internal logic error)

	// Service unavailable errors (503)
	CodeServiceUnavailable: http.StatusServiceUnavailable,
}

func GetHTTPStatus(code string) int {
	if status, ok := httpStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError // default to Internal Server Error
}
