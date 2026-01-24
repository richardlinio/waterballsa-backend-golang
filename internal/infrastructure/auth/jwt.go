package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

// NewJWTMiddleware creates and configures a new JWT middleware instance
// It accepts middleware functions via dependency injection for token operations
func NewJWTMiddleware(
	cfg config.JWTConfig,
	identityHandler func(*gin.Context) any,
	authorizer func(*gin.Context, any) bool,
	unauthorized func(*gin.Context, int, string),
) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       cfg.Realm,
		Key:         cfg.Secret,
		Timeout:     cfg.AccessTokenTimeout,
		MaxRefresh:  cfg.RefreshTokenTimeout,
		IdentityKey: "JWT_PAYLOAD",

		// Token lookup configuration
		// Supports both Authorization header and cookie for flexibility
		TokenLookup:   fmt.Sprintf("header: Authorization, cookie: %s", cfg.AccessTokenCookieName),
		TokenHeadName: "Bearer",

		// Cookie configuration for access token
		SendCookie:     true,
		CookieMaxAge:   cfg.AccessTokenTimeout,
		CookieHTTPOnly: true,
		SecureCookie:   cfg.SecureCookie, // HTTPS only in production
		CookieSameSite: http.SameSiteLaxMode,
		CookieDomain:   cfg.CookieDomain,
		CookieName:     cfg.AccessTokenCookieName,

		// Refresh token cookie configuration
		RefreshTokenCookieName: cfg.RefreshTokenCookieName,

		// Use injected middleware functions for token operations
		IdentityHandler: identityHandler,
		Authorizer:      authorizer,
		Unauthorized:    unauthorized,

		// HTTPStatusMessageFunc provides custom error messages
		HTTPStatusMessageFunc: func(c *gin.Context, e error) string {
			return e.Error()
		},

		TimeFunc: time.Now,
	})
}
