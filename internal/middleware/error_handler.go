package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Process request first
		c.Next()

		// Step 2: Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Step 3: Check if response has already been written (avoid duplicate writes)
		if c.Writer.Written() {
			return
		}

		// Step 4: Get the last error
		err := c.Errors.Last().Err

		// Step 5: Determine error type and respond
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			// Custom error - log and return structured error
			logger.Warn("Application error",
				"code", appErr.Code,
				"message", appErr.Message,
				"status", appErr.HTTPStatus,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			// If there's an underlying error, log it at error level
			if appErr.Err != nil {
				logger.Error("Underlying error", "error", appErr.Err)
			}

			response := dto.ErrorResponse{
				Code:  appErr.Code,
				Error: appErr.Message,
			}

			// Parse validation error details
			if appErr.Code == apperror.CodeValidationFailed && appErr.Err != nil {
				response.Details = parseValidationErrors(appErr.Err)
			}

			c.JSON(appErr.HTTPStatus, response)
		} else {
			// Unexpected error - log full error and return generic error
			logger.Error("Unexpected error",
				"error", err,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:  apperror.CodeInternalError,
				Error: apperror.GetMessage(apperror.CodeInternalError),
			})
		}
	}
}

// parseValidationErrors parses Gin validation errors
func parseValidationErrors(err error) map[string]string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make(map[string]string)
		for _, fe := range ve {
			details[fe.Field()] = getValidationMessage(fe)
		}
		return details
	}
	return nil
}

// getValidationMessage returns Chinese error messages based on validation rules
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "此欄位為必填"
	case "min":
		return fmt.Sprintf("長度至少需要 %s 個字元", fe.Param())
	case "max":
		return fmt.Sprintf("長度不可超過 %s 個字元", fe.Param())
	default:
		return "格式無效"
	}
}
