package dto

import "github.com/richardlinio/waterballsa-backend-golang/internal/model"

// UserProfileResponse represents the response for GET /users/me
type UserProfileResponse struct {
	ID               int64  `json:"id"`
	Username         string `json:"username"`
	ExperiencePoints int32  `json:"experiencePoints"`
	Level            int32  `json:"level"`
	Role             string `json:"role"`
}

// ToUserProfileResponse converts User domain model to UserProfileResponse DTO
func ToUserProfileResponse(user *model.User) UserProfileResponse {
	return UserProfileResponse{
		ID:               user.ID,
		Username:         user.Username,
		ExperiencePoints: user.ExperiencePoints,
		Level:            user.Level,
		Role:             user.Role,
	}
}
