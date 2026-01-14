package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

type AuthHandler struct {
	authService    service.AuthService
	tokenService   *auth.TokenService
	jwtMiddleware  *jwt.GinJWTMiddleware
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewAuthHandler(
	authService service.AuthService,
	tokenService *auth.TokenService,
	jwtMiddleware *jwt.GinJWTMiddleware,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		tokenService:   tokenService,
		jwtMiddleware:  jwtMiddleware,
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

	h.tokenService.SetCookie(c, result.Token, result.Expire)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: result.Token,
		User:        result.UserInfo,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := h.tokenService.ExtractToken(c)
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

	h.tokenService.ClearCookie(c)

	c.JSON(http.StatusOK, dto.LogoutResponse{
		Message: "登出成功",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	h.jwtMiddleware.RefreshHandler(c)
}
