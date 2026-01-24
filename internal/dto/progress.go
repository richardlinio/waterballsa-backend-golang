package dto

import "github.com/linporu/waterballsa-backend-golang/internal/model"

// UpdateProgressRequest for PUT /users/{userId}/missions/{missionId}/progress
type UpdateProgressRequest struct {
	WatchPositionSeconds int `json:"watchPositionSeconds" binding:"required,min=0"`
}

// UserMissionProgressResponse for progress endpoints
type UserMissionProgressResponse struct {
	MissionID            int64  `json:"missionId"`
	Status               string `json:"status"`
	WatchPositionSeconds int    `json:"watchPositionSeconds"`
}

// ToUserMissionProgressResponse converts domain model to DTO
func ToUserMissionProgressResponse(progress *model.UserMissionProgress) UserMissionProgressResponse {
	return UserMissionProgressResponse{
		MissionID:            progress.MissionID,
		Status:               progress.Status,
		WatchPositionSeconds: progress.WatchPositionSeconds,
	}
}

// DeliverResponse for POST /users/{userId}/missions/{missionId}/progress/deliver
type DeliverResponse struct {
	Message          string `json:"message"`
	ExperienceGained int    `json:"experienceGained"`
	TotalExperience  int32  `json:"totalExperience"`
	CurrentLevel     int32  `json:"currentLevel"`
}

// ToDeliverResponse converts domain model to DTO
func ToDeliverResponse(result *model.MissionDeliveryResult) DeliverResponse {
	return DeliverResponse{
		Message:          result.Message,
		ExperienceGained: result.ExperienceGained,
		TotalExperience:  result.TotalExperience,
		CurrentLevel:     result.CurrentLevel,
	}
}
