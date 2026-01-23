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

// JourneyDetailResponse is the response for GET /journeys/{journeyId} endpoint
type JourneyDetailResponse struct {
	ID            int64        `json:"id"`
	Slug          string       `json:"slug"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	CoverImageURL string       `json:"coverImageUrl"`
	TeacherName   string       `json:"teacherName"`
	Price         float64      `json:"price"`
	Chapters      []ChapterDTO `json:"chapters"`
}

// ChapterDTO represents a chapter with its missions
type ChapterDTO struct {
	ID         int64               `json:"id"`
	Title      string              `json:"title"`
	OrderIndex int                 `json:"orderIndex"`
	Missions   []MissionSummaryDTO `json:"missions"`
}

// MissionSummaryDTO represents a mission summary in a chapter
type MissionSummaryDTO struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	AccessLevel string `json:"accessLevel"`
	OrderIndex  int    `json:"orderIndex"`
}
