package apperror

import "net/http"

var httpStatusMap = map[string]int{
	// Validation errors (400 - Bad Request)
	CodeValidationFailed:     http.StatusBadRequest,
	CodePasswordTooLong:      http.StatusBadRequest,
	CodeInvalidWatchPosition: http.StatusBadRequest,
	CodeMissionNotCompleted:  http.StatusBadRequest,

	// Auth errors (401 - Unauthorized)
	CodeAuthFailed:   http.StatusUnauthorized, // Login failures, wrong credentials, rate limited
	CodeUnauthorized: http.StatusUnauthorized, // Token invalid/missing

	// Forbidden errors (403 - Forbidden)
	CodeProgressForbidden: http.StatusForbidden,

	// Not found errors (404 - Not Found)
	CodeJourneyNotFound:  http.StatusNotFound,
	CodeMissionNotFound:  http.StatusNotFound,
	CodeProgressNotFound: http.StatusNotFound,

	// Conflict errors (409 - Conflict)
	CodeUsernameExists:          http.StatusConflict,
	CodeMissionAlreadyDelivered: http.StatusConflict,

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
