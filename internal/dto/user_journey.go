package dto

// UserJourneyItem represents a single purchased journey
type UserJourneyItem struct {
	JourneyID     int64  `json:"journeyId"`
	JourneyTitle  string `json:"journeyTitle"`
	JourneySlug   string `json:"journeySlug"`
	CoverImageURL string `json:"coverImageUrl"`
	TeacherName   string `json:"teacherName"`
	PurchasedAt   int64  `json:"purchasedAt"` // Unix milliseconds
	OrderNumber   string `json:"orderNumber"`
}

// UserJourneyListResponse is the response for GET /users/:userId/journeys
type UserJourneyListResponse struct {
	Journeys []UserJourneyItem `json:"journeys"`
}

func ToUserJourneyListResponse(journeys []UserJourneyItem) UserJourneyListResponse {
	return UserJourneyListResponse{
		Journeys: journeys,
	}
}
