package auth

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// TokenService handles all JWT token-related operations
type TokenService struct {
	config config.JWTConfig
}

// NewTokenService creates a new TokenService instance
func NewTokenService(config config.JWTConfig) *TokenService {
	return &TokenService{config: config}
}

// Generate creates a new JWT token for the user
func (ts *TokenService) Generate(user *model.User) (string, time.Time, error) {
	expireTime := time.Now().Add(ts.config.AccessTokenTimeout)

	claims := gojwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"experience": user.ExperiencePoints,
		"exp":        expireTime.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ts.config.Secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expireTime, nil
}

// SetCookie sets the JWT token as an HTTP-only cookie
func (ts *TokenService) SetCookie(c *gin.Context, token string, expire time.Time) {
	maxAge := int(time.Until(expire).Seconds())

	c.SetCookie(
		"jwt",                  // name
		token,                  // value
		maxAge,                 // maxAge
		"/",                    // path
		ts.config.CookieDomain, // domain
		ts.config.SecureCookie, // secure
		true,                   // httpOnly
	)
}

// ClearCookie removes the JWT cookie
func (ts *TokenService) ClearCookie(c *gin.Context) {
	c.SetCookie(
		"jwt",
		"",
		-1,
		"/",
		ts.config.CookieDomain,
		ts.config.SecureCookie,
		true,
	)
}

// ExtractToken gets token from Authorization header or cookie
func (ts *TokenService) ExtractToken(c *gin.Context) (string, error) {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Remove "Bearer " prefix
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:], nil
		}
	}

	// Try cookie
	token, err := c.Cookie("jwt")
	if err != nil {
		return "", fmt.Errorf("token not found in header or cookie: %w", err)
	}

	return token, nil
}

// ExtractIdentity parses JWT token and extracts user identity (for middleware)
func (ts *TokenService) ExtractIdentity(c *gin.Context) interface{} {
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
func (ts *TokenService) Authorize(_ *gin.Context, data interface{}) bool {
	if _, ok := data.(*model.User); ok {
		return true
	}
	return false
}

// HandleUnauthorized handles unauthorized access (for middleware)
func (ts *TokenService) HandleUnauthorized(c *gin.Context, _ int, _ string) {
	_ = c.Error(apperror.Unauthorized())
}
