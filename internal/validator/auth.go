package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9@$!%*?&#]+$`)
)

func RegisterAuthValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("alphanum_underscore", validateUsernameCharset); err != nil {
			return err
		}
		if err := v.RegisterValidation("password_charset", validatePasswordCharset); err != nil {
			return err
		}
	}
	return nil
}

func validateUsernameCharset(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegex.MatchString(username)
}

func validatePasswordCharset(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return passwordRegex.MatchString(password)
}
