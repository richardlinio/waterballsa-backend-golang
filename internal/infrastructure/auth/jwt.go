package auth

import (
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
	identityHandler func(*gin.Context) interface{},
	authorizer func(*gin.Context, interface{}) bool,
	unauthorized func(*gin.Context, int, string),
) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "waterballsa",
		Key:         cfg.Secret,
		Timeout:     cfg.AccessTokenTimeout,
		MaxRefresh:  cfg.RefreshTokenTimeout,
		IdentityKey: "user_id",

		// Token lookup configuration
		// Supports both Authorization header and cookie for flexibility
		TokenLookup:   "header: Authorization, cookie: jwt",
		TokenHeadName: "Bearer",

		// Cookie configuration for access token
		SendCookie:     true,
		CookieMaxAge:   cfg.AccessTokenTimeout,
		CookieHTTPOnly: true,
		SecureCookie:   cfg.SecureCookie, // HTTPS only in production
		CookieSameSite: http.SameSiteLaxMode,
		CookieDomain:   cfg.CookieDomain,
		CookieName:     "jwt",

		// Refresh token cookie configuration
		RefreshTokenCookieName: "refresh_token",

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
