package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

// NewJWTMiddleware creates and configures a new JWT middleware instance
func NewJWTMiddleware(
	cfg config.JWTConfig,
	authService service.AuthService,
	logger *slog.Logger,
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

		// Authenticator validates user credentials during login
		Authenticator: func(c *gin.Context) (any, error) {
			var req dto.LoginRequest

			// Bind and validate JSON request
			if err := c.ShouldBindJSON(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			// Create context with timeout for service call
			ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
			defer cancel()

			// Authenticate user via service layer
			user, err := authService.Login(ctx, req)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// Store user data in context for LoginResponse to access
			c.Set("authenticated_user", user)

			return user, nil
		},

		// PayloadFunc sets JWT claims from authenticated user data
		PayloadFunc: func(data interface{}) gojwt.MapClaims {
			if user, ok := data.(*model.User); ok {
				return gojwt.MapClaims{
					"user_id":    user.ID,
					"username":   user.Username,
					"role":       user.Role,
					"experience": user.ExperiencePoints,
				}
			}
			return gojwt.MapClaims{}
		},

		// IdentityHandler extracts user identity from JWT claims
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			// Extract user_id from claims
			userID, ok := claims["user_id"].(float64)
			if !ok {
				return nil
			}

			// Extract username from claims
			username, ok := claims["username"].(string)
			if !ok {
				return nil
			}

			// Extract role from claims
			role, ok := claims["role"].(string)
			if !ok {
				return nil
			}

			return &model.User{
				ID:       int64(userID),
				Username: username,
				Role:     role,
			}
		},

		// Authorizer performs authorization check (default: allow all authenticated users)
		Authorizer: func(c *gin.Context, data interface{}) bool {
			// Basic authorization: allow all authenticated users
			// Can be extended with role-based checks if needed
			if _, ok := data.(*model.User); ok {
				return true
			}
			return false
		},

		// LoginResponse customizes the login success response
		LoginResponse: func(c *gin.Context, token *core.Token) {
			// Get user data from context (set by Authenticator)
			userData, exists := c.Get("authenticated_user")
			if !exists {
				_ = c.Error(apperror.AuthStateError(
					fmt.Errorf("authenticated user not found in context")))
				return
			}

			user, ok := userData.(*model.User)
			if !ok {
				_ = c.Error(apperror.AuthStateError(
					fmt.Errorf("invalid user data type in context: %T", userData)))
				return
			}

			c.JSON(http.StatusOK, dto.LoginResponse{
				AccessToken: token.AccessToken,
				User: dto.UserInfo{
					ID:         user.ID,
					Username:   user.Username,
					Experience: user.ExperiencePoints,
				},
			})
		},

		// RefreshResponse customizes the token refresh response
		RefreshResponse: func(c *gin.Context, token *core.Token) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token.AccessToken,
				"expire": time.Unix(token.ExpiresAt, 0).Format(time.RFC3339),
			})
		},

		// LogoutResponse customizes the logout response
		LogoutResponse: func(c *gin.Context) {
			c.JSON(http.StatusOK, dto.LogoutResponse{
				Message: "登出成功",
			})
		},

		// Unauthorized handles unauthorized access attempts
		Unauthorized: func(c *gin.Context, code int, message string) {
			// Check if this is a login failure (ErrFailedAuthentication)
			if message == jwt.ErrFailedAuthentication.Error() {
				_ = c.Error(apperror.AuthFailed())
			} else {
				_ = c.Error(apperror.Unauthorized())
			}
		},

		// HTTPStatusMessageFunc provides custom error messages
		HTTPStatusMessageFunc: func(c *gin.Context, e error) string {
			return e.Error()
		},

		TimeFunc: time.Now,
	})
}
