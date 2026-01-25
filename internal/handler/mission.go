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
	GetDetail(ctx context.Context, missionID int64) (*model.MissionDetail, error)
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

	// Get mission detail from service
	detail, err := h.missionService.GetDetail(ctx, missionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Check access level and authentication
	if detail.Mission.AccessLevel == "AUTHENTICATED" || detail.Mission.AccessLevel == "PURCHASED" {
		user := util.GetAuthenticatedUser(c)
		if user == nil {
			_ = c.Error(apperror.Unauthorized())
			return
		}

		// TODO: For PURCHASED missions, implement purchase verification logic
		// Need to check if user.ID has purchased the journey (detail.JourneyID)
		// Example implementation:
		//   if detail.Mission.AccessLevel == "PURCHASED" {
		//       hasPurchased, err := h.purchaseService.HasUserPurchasedJourney(ctx, user.ID, detail.JourneyID)
		//       if err != nil {
		//           _ = c.Error(apperror.DatabaseError(err))
		//           return
		//       }
		//       if !hasPurchased {
		//           _ = c.Error(apperror.Forbidden())
		//           return
		//       }
		//   }
	}

	response := dto.ToMissionDetailResponse(detail)
	c.JSON(http.StatusOK, response)
}
