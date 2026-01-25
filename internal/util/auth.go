package util

import (
	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
)

// GetAuthenticatedUser retrieves the authenticated user from the Gin context.
// Returns nil if the user is not found or the payload is invalid.
func GetAuthenticatedUser(c *gin.Context) *model.User {
	payload, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return nil
	}

	user, ok := payload.(*model.User)
	if !ok {
		return nil
	}

	return user
}
