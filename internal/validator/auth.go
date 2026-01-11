package validator

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9@$!%*?&#]+$`)

	// Ensure validators are registered only once across all calls
	registerOnce sync.Once
	registerErr  error
)

// RegisterAuthValidators registers custom validators for authentication fields.
// This function is idempotent and safe to call multiple times - the registration
// logic will only execute once regardless of how many times it's called.
func RegisterAuthValidators() error {
	registerOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			if err := v.RegisterValidation("alphanum_underscore", validateUsernameCharset); err != nil {
				registerErr = fmt.Errorf("failed to register alphanum_underscore validator: %w", err)
				return
			}
			if err := v.RegisterValidation("password_charset", validatePasswordCharset); err != nil {
				registerErr = fmt.Errorf("failed to register password_charset validator: %w", err)
				return
			}
		}
	})
	return registerErr
}

func validateUsernameCharset(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegex.MatchString(username)
}

func validatePasswordCharset(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return passwordRegex.MatchString(password)
}
