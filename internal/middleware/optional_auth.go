package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

// OptionalJWTAuth attempts to authenticate the request but allows it to proceed even if authentication fails
// If a valid JWT token is present, it sets the user identity in context (same as JWTAuth)
// If no valid token is present, the request continues as a guest (user identity is nil)
func OptionalJWTAuth(jwtMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to parse and validate token
		token, err := jwtMiddleware.ParseToken(c)
		if err != nil || token == nil || !token.Valid {
			// No valid token - continue as guest
			c.Next()
			return
		}

		// Valid token - extract claims and set in context
		claims := jwt.ExtractClaimsFromToken(token)
		if claims == nil {
			// No claims - continue as guest
			c.Next()
			return
		}

		// Store claims in context under "JWT_PAYLOAD" (required by jwt.ExtractClaims helper)
		c.Set("JWT_PAYLOAD", claims)

		// Extract identity using the configured IdentityHandler
		// This will call ExtractIdentity which uses jwt.ExtractClaims(c)
		identity := jwtMiddleware.IdentityHandler(c)
		if identity != nil {
			c.Set(jwtMiddleware.IdentityKey, identity)
		}

		c.Next()
	}
}
