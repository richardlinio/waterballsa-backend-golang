package middleware

import (
	"context"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// tokenBlacklistChecker is used by BlacklistChecker middleware to check token blacklist
type tokenBlacklistChecker interface {
	IsInvalidated(ctx context.Context, jti string) (bool, error)
}

// JWTAuth returns a middleware that validates JWT tokens
// and sets user identity in the request context.
func JWTAuth(jwtMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return jwtMiddleware.MiddlewareFunc()
}

// ExtractIdentity parses JWT token and extracts user identity (for middleware)
func ExtractIdentity(c *gin.Context) any {
	claims := jwt.ExtractClaims(c) // gin-jwt helper

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil
	}

	return &model.User{
		ID:       int64(userID),
		Username: username,
		Role:     role,
	}
}

// Authorize checks if user is authorized (for middleware)
func Authorize(_ *gin.Context, data any) bool {
	if _, ok := data.(*model.User); ok {
		return true
	}
	return false
}

// HandleUnauthorized handles unauthorized access (for middleware)
func HandleUnauthorized(c *gin.Context, _ int, _ string) {
	_ = c.Error(apperror.Unauthorized())
}

// BlacklistChecker returns a middleware that checks if the token's JTI is blacklisted
func BlacklistChecker(checker tokenBlacklistChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)

		jti, ok := claims["jti"].(string)
		if !ok {
			// Token doesn't have JTI claim (legacy token), allow through
			c.Next()
			return
		}

		invalidated, err := checker.IsInvalidated(c.Request.Context(), jti)
		if err != nil {
			_ = c.Error(apperror.InternalServerError(err))
			c.Abort()
			return
		}

		if invalidated {
			_ = c.Error(apperror.Unauthorized())
			c.Abort()
			return
		}

		c.Next()
	}
}
