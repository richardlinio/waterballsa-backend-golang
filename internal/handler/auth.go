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
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

type AuthHandler struct {
	authService    service.AuthService
	jwtMiddleware  *jwt.GinJWTMiddleware
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewAuthHandler(
	authService service.AuthService,
	jwtMiddleware *jwt.GinJWTMiddleware,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		jwtMiddleware:  jwtMiddleware,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Call service to register user
	userID, err := h.authService.Register(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "註冊成功",
		UserID:  userID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.jwtMiddleware.LoginHandler(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.jwtMiddleware.LogoutHandler(c)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	h.jwtMiddleware.RefreshHandler(c)
}
