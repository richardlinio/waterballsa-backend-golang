package apperror

type AppError struct {
	Code       string // Error code (e.g., "ERR_USERNAME_EXISTS")
	Message    string // Error message
	HTTPStatus int    // HTTP status code
	Err        error  // Underlying error (for logging)
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with HTTP status code automatically mapped from code
func New(code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    GetMessage(code),
		HTTPStatus: GetHTTPStatus(code),
	}
}

// NewWithError creates a new AppError with an underlying error
func NewWithError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    GetMessage(code),
		HTTPStatus: GetHTTPStatus(code),
		Err:        err,
	}
}

// Predefined error constructor functions - no need to specify HTTP status code

func ValidationFailed() *AppError {
	return New(CodeValidationFailed)
}

func UsernameExists() *AppError {
	return New(CodeUsernameExists)
}

func PasswordTooLong() *AppError {
	return New(CodePasswordTooLong)
}

// AuthFailed returns a generic authentication failure error
// Use this for: login failures, wrong username/password, rate limiting
// to avoid leaking information about which part of authentication failed
func AuthFailed() *AppError {
	return New(CodeAuthFailed)
}

// Unauthorized returns an error for invalid/missing tokens
// Use this for: token validation failures in protected endpoints
func Unauthorized() *AppError {
	return New(CodeUnauthorized)
}

// AuthStateError returns an error for authentication state errors
// Use this for: internal logic errors in auth flow (e.g., missing user data in context)
func AuthStateError(err error) *AppError {
	return NewWithError(CodeAuthStateError, err)
}

func InternalServerError(err error) *AppError {
	return NewWithError(CodeInternalServerError, err)
}

func DatabaseError(err error) *AppError {
	return NewWithError(CodeDatabaseError, err)
}

func ServiceUnavailable() *AppError {
	return New(CodeServiceUnavailable)
}

func RateLimitExceeded() *AppError {
	return New(CodeRateLimitExceeded)
}

func JourneyNotFound() *AppError {
	return New(CodeJourneyNotFound)
}
