package dto

import "github.com/richardlinio/waterballsa-backend-golang/internal/model"

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

// ToJourneyListResponse converts journey domain models to JourneyListResponse DTO
func ToJourneyListResponse(journeys []*model.Journey) JourneyListResponse {
	if journeys == nil {
		return JourneyListResponse{
			Journeys: []JourneyListItem{},
		}
	}
	items := make([]JourneyListItem, 0, len(journeys))
	for _, journey := range journeys {
		items = append(items, ToJourneyListItem(journey))
	}
	return JourneyListResponse{
		Journeys: items,
	}
}

// ToJourneyListItem converts a single journey domain model to JourneyListItem DTO
func ToJourneyListItem(journey *model.Journey) JourneyListItem {
	return JourneyListItem{
		ID:            journey.ID,
		Slug:          journey.Slug,
		Title:         journey.Title,
		Description:   journey.Description,
		CoverImageURL: journey.CoverImageURL,
		TeacherName:   journey.TeacherName,
		Price:         journey.Price,
	}
}

// ToJourneyDetailResponse converts JourneyDetail domain model to JourneyDetailResponse DTO
func ToJourneyDetailResponse(detail *model.JourneyDetail) JourneyDetailResponse {
	return JourneyDetailResponse{
		ID:            detail.Journey.ID,
		Slug:          detail.Journey.Slug,
		Title:         detail.Journey.Title,
		Description:   detail.Journey.Description,
		CoverImageURL: detail.Journey.CoverImageURL,
		TeacherName:   detail.Journey.TeacherName,
		Price:         detail.Journey.Price,
		Chapters:      ToChapterDTOs(detail.Chapters, detail.MissionsByChapter),
	}
}

// ToChapterDTOs converts chapters and missions to ChapterDTO slice
func ToChapterDTOs(chapters []*model.Chapter, missionsByChapter map[int64][]*model.Mission) []ChapterDTO {
	if chapters == nil {
		return []ChapterDTO{}
	}
	chapterDTOs := make([]ChapterDTO, len(chapters))
	for i, chapter := range chapters {
		chapterDTOs[i] = ToChapterDTO(chapter, missionsByChapter[chapter.ID])
	}
	return chapterDTOs
}

// ToChapterDTO converts a single chapter with missions to ChapterDTO
func ToChapterDTO(chapter *model.Chapter, missions []*model.Mission) ChapterDTO {
	return ChapterDTO{
		ID:         chapter.ID,
		Title:      chapter.Title,
		OrderIndex: chapter.OrderIndex,
		Missions:   ToMissionSummaryDTOs(missions),
	}
}

// ToMissionSummaryDTOs converts missions to MissionSummaryDTO slice
func ToMissionSummaryDTOs(missions []*model.Mission) []MissionSummaryDTO {
	if missions == nil {
		return []MissionSummaryDTO{}
	}
	missionDTOs := make([]MissionSummaryDTO, len(missions))
	for i, mission := range missions {
		missionDTOs[i] = ToMissionSummaryDTO(mission)
	}
	return missionDTOs
}

// ToMissionSummaryDTO converts a single mission to MissionSummaryDTO
func ToMissionSummaryDTO(mission *model.Mission) MissionSummaryDTO {
	return MissionSummaryDTO{
		ID:          mission.ID,
		Type:        mission.Type,
		Title:       mission.Title,
		AccessLevel: mission.AccessLevel,
		OrderIndex:  mission.OrderIndex,
	}
}
