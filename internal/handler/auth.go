package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
	refreshCookiePath      = "/auth/refresh"
	bearerPrefix           = "Bearer "
)

type AuthHandler struct {
	authService    service.AuthService
	jwtMiddleware  *jwt.GinJWTMiddleware
	jwtConfig      config.JWTConfig
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewAuthHandler(
	authService service.AuthService,
	jwtMiddleware *jwt.GinJWTMiddleware,
	jwtConfig config.JWTConfig,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		jwtMiddleware:  jwtMiddleware,
		jwtConfig:      jwtConfig,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	userID, err := h.authService.Register(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "註冊成功",
		UserID:  userID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	result, err := h.authService.Login(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Set access token cookie
	h.setAccessTokenCookie(c, result.AccessToken, result.AccessTokenExpire)

	// Set refresh token cookie (restricted path)
	h.setRefreshTokenCookie(c, result.RefreshToken, result.RefreshTokenExpire)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: result.AccessToken,
		User:        result.UserInfo,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract access token (may fail if already expired, that's ok)
	accessToken, _ := h.extractAccessToken(c)

	// Extract refresh token
	refreshToken, _ := h.extractRefreshToken(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	if err := h.authService.Logout(ctx, accessToken, refreshToken); err != nil {
		_ = c.Error(err)
		return
	}

	// Clear both cookies
	h.clearAccessTokenCookie(c)
	h.clearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, dto.LogoutResponse{
		Message: "登出成功",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := h.extractRefreshToken(c)
	if err != nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	result, err := h.authService.Refresh(ctx, refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Set new access token cookie
	h.setAccessTokenCookie(c, result.AccessToken, result.AccessTokenExpire)

	// Set new refresh token cookie (token rotation)
	h.setRefreshTokenCookie(c, result.RefreshToken, result.RefreshTokenExpire)

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken: result.AccessToken,
		User:        result.UserInfo,
	})
}

// setAccessTokenCookie sets the access token as a cookie (not HTTP-only for frontend access)
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string, expire time.Time) {
	maxAge := int(time.Until(expire).Seconds())

	c.SetCookie(
		accessTokenCookieName,
		token,
		maxAge,
		"/",
		h.jwtConfig.CookieDomain,
		h.jwtConfig.SecureCookie,
		false, // httpOnly = false, allows frontend to read for login state checking
	)
}

// setRefreshTokenCookie sets the refresh token as an HTTP-only cookie with restricted path
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string, expire time.Time) {
	maxAge := int(time.Until(expire).Seconds())

	c.SetCookie(
		refreshTokenCookieName,
		token,
		maxAge,
		refreshCookiePath, // Restricted to /auth/refresh path only
		h.jwtConfig.CookieDomain,
		h.jwtConfig.SecureCookie,
		true, // httpOnly
	)
}

// clearAccessTokenCookie removes the access token cookie
func (h *AuthHandler) clearAccessTokenCookie(c *gin.Context) {
	c.SetCookie(
		accessTokenCookieName,
		"",
		-1,
		"/",
		h.jwtConfig.CookieDomain,
		h.jwtConfig.SecureCookie,
		false, // Must match httpOnly setting when clearing
	)
}

// clearRefreshTokenCookie removes the refresh token cookie
func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		refreshCookiePath,
		h.jwtConfig.CookieDomain,
		h.jwtConfig.SecureCookie,
		true,
	)
}

// extractAccessToken gets access token from Authorization header or cookie
func (h *AuthHandler) extractAccessToken(c *gin.Context) (string, error) {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if token := strings.TrimPrefix(authHeader, bearerPrefix); token != authHeader {
		return token, nil
	}

	// Try cookie
	return c.Cookie(accessTokenCookieName)
}

// extractRefreshToken gets refresh token from cookie only
func (h *AuthHandler) extractRefreshToken(c *gin.Context) (string, error) {
	return c.Cookie(refreshTokenCookieName)
}
