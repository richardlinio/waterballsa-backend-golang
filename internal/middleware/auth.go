package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

// JWTAuth returns a middleware that validates JWT tokens
// and sets user identity in the request context.
func JWTAuth(jwtMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return jwtMiddleware.MiddlewareFunc()
}
