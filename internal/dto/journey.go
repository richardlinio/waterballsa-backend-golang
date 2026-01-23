package dto

// JourneyListResponse is the response for GET /journeys endpoint
type JourneyListResponse struct {
	Journeys []JourneyListItem `json:"journeys"`
}

// JourneyListItem represents a single journey in the list
type JourneyListItem struct {
	ID            int64   `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	CoverImageURL string  `json:"coverImageUrl"`
	TeacherName   string  `json:"teacherName"`
	Price         float64 `json:"price"`
}
