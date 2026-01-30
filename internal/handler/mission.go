package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/service"
	"github.com/richardlinio/waterballsa-backend-golang/internal/util"
)

// missionService defines the mission service operations needed by the handler
type missionService interface {
	GetDetail(ctx context.Context, missionID int64, userID *int64) (*model.MissionDetail, error)
}

type MissionHandler struct {
	missionService missionService
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewMissionHandler(
	missionService *service.MissionService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *MissionHandler {
	return &MissionHandler{
		missionService: missionService,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *MissionHandler) GetMissionDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Parse missionId from path parameter
	missionIDStr := c.Param("missionId")
	missionID, err := strconv.ParseInt(missionIDStr, 10, 64)
	if err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	// Get authenticated user (may be nil for guest users)
	user := util.GetAuthenticatedUser(c)
	var userID *int64
	if user != nil {
		userID = &user.ID
	}

	// Get mission detail from service (service handles access control)
	detail, err := h.missionService.GetDetail(ctx, missionID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToMissionDetailResponse(detail)
	c.JSON(http.StatusOK, response)
}
