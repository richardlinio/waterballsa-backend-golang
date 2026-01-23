package dto

import (
	"fmt"
	"strings"

	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// MissionDetailResponse is the response for GET /journeys/{journeyId}/missions/{missionId} endpoint
type MissionDetailResponse struct {
	ID          int64                `json:"id"`
	ChapterID   int64                `json:"chapterId"`
	JourneyID   int64                `json:"journeyId"`
	Type        string               `json:"type"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	AccessLevel string               `json:"accessLevel"`
	CreatedAt   int64                `json:"createdAt"`
	VideoLength string               `json:"videoLength,omitempty"`
	Reward      MissionRewardDTO     `json:"reward"`
	Resource    []MissionResourceDTO `json:"resource"`
}

// MissionRewardDTO represents the reward information in mission detail
type MissionRewardDTO struct {
	Exp int `json:"exp"`
}

// MissionResourceDTO represents a single resource in mission detail
type MissionResourceDTO struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	ResourceURL     string `json:"resourceUrl,omitempty"`
	ResourceContent string `json:"resourceContent,omitempty"`
	DurationSeconds *int   `json:"durationSeconds,omitempty"`
}

// ToMissionDetailResponse converts MissionDetail domain model to MissionDetailResponse DTO
func ToMissionDetailResponse(detail *model.MissionDetail) MissionDetailResponse {
	response := MissionDetailResponse{
		ID:          detail.Mission.ID,
		ChapterID:   detail.Mission.ChapterID,
		JourneyID:   detail.JourneyID,
		Type:        detail.Mission.Type,
		Title:       detail.Mission.Title,
		Description: detail.Mission.Description,
		AccessLevel: detail.Mission.AccessLevel,
		CreatedAt:   detail.Mission.CreatedAt.UnixMilli(),
		Reward:      ToMissionRewardDTO(detail.Reward),
		Resource:    ToMissionResourceDTOs(detail.Resources),
	}

	// Calculate video length if mission type is VIDEO
	if detail.Mission.Type == "VIDEO" {
		response.VideoLength = calculateVideoLength(detail.Resources)
	}

	return response
}

// ToMissionRewardDTO converts a reward domain model to MissionRewardDTO
func ToMissionRewardDTO(reward *model.Reward) MissionRewardDTO {
	if reward == nil {
		return MissionRewardDTO{
			Exp: 0,
		}
	}

	return MissionRewardDTO{
		Exp: reward.RewardValue,
	}
}

// ToMissionResourceDTOs converts mission resources to MissionResourceDTO slice
func ToMissionResourceDTOs(resources []*model.MissionResource) []MissionResourceDTO {
	if resources == nil {
		return []MissionResourceDTO{}
	}

	resourceDTOs := make([]MissionResourceDTO, len(resources))
	for i, resource := range resources {
		resourceDTOs[i] = ToMissionResourceDTO(resource)
	}

	return resourceDTOs
}

// ToMissionResourceDTO converts a single mission resource to MissionResourceDTO
func ToMissionResourceDTO(resource *model.MissionResource) MissionResourceDTO {
	return MissionResourceDTO{
		ID:              int(resource.ID),
		Type:            strings.ToLower(resource.ResourceType), // Convert VIDEO to video, ARTICLE to article, FORM to form
		ResourceURL:     resource.ResourceURL,
		ResourceContent: resource.ResourceContent,
		DurationSeconds: resource.DurationSeconds,
	}
}

// calculateVideoLength calculates the total video length from resources
// Returns formatted string in "MM:SS" format
func calculateVideoLength(resources []*model.MissionResource) string {
	totalSeconds := 0

	for _, resource := range resources {
		if resource.ResourceType == "VIDEO" && resource.DurationSeconds != nil {
			totalSeconds += *resource.DurationSeconds
		}
	}

	if totalSeconds == 0 {
		return ""
	}

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
