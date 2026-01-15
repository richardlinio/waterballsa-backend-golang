package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

type AuthHandler struct {
	authService    service.AuthService
	tokenGenerator auth.TokenGenerator
	jwtMiddleware  *jwt.GinJWTMiddleware
	jwtConfig      config.JWTConfig
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewAuthHandler(
	authService service.AuthService,
	tokenGenerator auth.TokenGenerator,
	jwtMiddleware *jwt.GinJWTMiddleware,
	jwtConfig config.JWTConfig,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		tokenGenerator: tokenGenerator,
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

	h.setCookie(c, result.Token, result.Expire)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: result.Token,
		User:        result.UserInfo,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := h.extractToken(c)
	if err != nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	if err := h.authService.Logout(ctx, token); err != nil {
		_ = c.Error(err)
		return
	}

	h.clearCookie(c)

	c.JSON(http.StatusOK, dto.LogoutResponse{
		Message: "登出成功",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	token, err := h.extractToken(c)
	if err != nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	result, err := h.authService.Refresh(ctx, token)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.setCookie(c, result.Token, result.Expire)

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken: result.Token,
		User:        result.UserInfo,
	})
}

// setCookie sets the JWT token as an HTTP-only cookie
func (h *AuthHandler) setCookie(c *gin.Context, token string, expire time.Time) {
	maxAge := int(time.Until(expire).Seconds())

	c.SetCookie(
		"jwt",                    // name
		token,                    // value
		maxAge,                   // maxAge
		"/",                      // path
		h.jwtConfig.CookieDomain, // domain
		h.jwtConfig.SecureCookie, // secure
		true,                     // httpOnly
	)
}

// clearCookie removes the JWT cookie
func (h *AuthHandler) clearCookie(c *gin.Context) {
	c.SetCookie(
		"jwt",
		"",
		-1,
		"/",
		h.jwtConfig.CookieDomain,
		h.jwtConfig.SecureCookie,
		true,
	)
}

// extractToken gets token from Authorization header or cookie
func (h *AuthHandler) extractToken(c *gin.Context) (string, error) {
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
