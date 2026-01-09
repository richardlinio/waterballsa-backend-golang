package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/requestid"
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

		// Step 5: Get request ID for logging
		reqID := requestid.Get(c)

		// Step 6: Determine error type and respond
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			// Custom error - log and return structured error
			logger.Warn("Application error",
				"request_id", reqID,
				"code", appErr.Code,
				"message", appErr.Message,
				"status", appErr.HTTPStatus,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"referer", c.Request.Referer(),
			)

			// If there's an underlying error, log it at error level
			if appErr.Err != nil {
				logger.Error("Underlying error",
					"request_id", reqID,
					"error", appErr.Err,
				)
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
				"request_id", reqID,
				"error", err,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"referer", c.Request.Referer(),
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
