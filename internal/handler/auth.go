package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

type AuthHandler struct {
	authService    *service.AuthService
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger, requestTimeout time.Duration) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "資料驗證失敗,請檢查輸入內容",
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Call service to register user
	userID, err := h.authService.Register(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "使用者名稱已存在",
			})
			return
		}

		h.logger.Error("Registration failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "註冊失敗,請稍後再試",
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "註冊成功",
		UserID:  userID,
	})
}
