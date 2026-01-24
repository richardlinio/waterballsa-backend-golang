package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

type progressService interface {
	GetProgress(ctx context.Context, userID, missionID int64) (*model.UserMissionProgress, error)
	UpdateProgress(ctx context.Context, userID, missionID int64, watchPositionSeconds int) (*model.UserMissionProgress, error)
	DeliverMission(ctx context.Context, userID, missionID int64) (*model.MissionDeliveryResult, error)
}

type ProgressHandler struct {
	progressService progressService
	logger          *slog.Logger
	requestTimeout  time.Duration
}

func NewProgressHandler(
	progressService *service.ProgressService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
		requestTimeout:  requestTimeout,
	}
}

func (h *ProgressHandler) GetProgress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Parse path parameters
	userID, missionID, err := h.parsePathParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Check authorization: user can only access their own progress
	authenticatedUser := h.getAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	if authenticatedUser.ID != userID {
		_ = c.Error(apperror.ProgressForbidden())
		return
	}

	// Get progress from service
	progress, err := h.progressService.GetProgress(ctx, userID, missionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToUserMissionProgressResponse(progress)
	c.JSON(http.StatusOK, response)
}

func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Parse path parameters
	userID, missionID, err := h.parsePathParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Check authorization: user can only update their own progress
	authenticatedUser := h.getAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	if authenticatedUser.ID != userID {
		_ = c.Error(apperror.ProgressForbidden())
		return
	}

	// Parse request body
	var req dto.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	// Update progress
	progress, err := h.progressService.UpdateProgress(ctx, userID, missionID, req.WatchPositionSeconds)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToUserMissionProgressResponse(progress)
	c.JSON(http.StatusOK, response)
}

func (h *ProgressHandler) DeliverMission(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Parse path parameters
	userID, missionID, err := h.parsePathParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Check authorization: user can only deliver their own missions
	authenticatedUser := h.getAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	if authenticatedUser.ID != userID {
		_ = c.Error(apperror.ProgressForbidden())
		return
	}

	// Deliver mission
	result, err := h.progressService.DeliverMission(ctx, userID, missionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToDeliverResponse(result)
	c.JSON(http.StatusOK, response)
}

// parsePathParams extracts and validates userId and missionId from path
func (h *ProgressHandler) parsePathParams(c *gin.Context) (userID int64, missionID int64, err error) {
	userIDStr := c.Param("userId")
	userID, err = strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, 0, apperror.NewWithError(apperror.CodeValidationFailed, err)
	}

	missionIDStr := c.Param("missionId")
	missionID, err = strconv.ParseInt(missionIDStr, 10, 64)
	if err != nil {
		return 0, 0, apperror.NewWithError(apperror.CodeValidationFailed, err)
	}

	return userID, missionID, nil
}

// getAuthenticatedUser extracts authenticated user from context
func (h *ProgressHandler) getAuthenticatedUser(c *gin.Context) *model.User {
	payload, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return nil
	}

	user, ok := payload.(*model.User)
	if !ok {
		return nil
	}

	return user
}
