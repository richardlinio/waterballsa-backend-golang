package service

import (
	"context"
	"errors"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

// journeyRepository provides journey data access operations for JourneyService
type journeyRepository interface {
	List(ctx context.Context) ([]*model.Journey, error)
	GetByID(ctx context.Context, id int64) (*model.Journey, error)
	ListChaptersByJourneyID(ctx context.Context, journeyID int64) ([]*model.Chapter, error)
	ListMissionsByChapterIDs(ctx context.Context, chapterIDs []int64) ([]*model.Mission, error)
}

// JourneyService implements journey business logic operations.
type JourneyService struct {
	journeyRepository journeyRepository
}

// NewJourneyService creates a new JourneyService instance.
func NewJourneyService(journeyRepository *repository.JourneyRepository) *JourneyService {
	return &JourneyService{
		journeyRepository: journeyRepository,
	}
}

// List retrieves all available journeys
func (s *JourneyService) List(ctx context.Context) ([]*model.Journey, error) {
	journeys, err := s.journeyRepository.List(ctx)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return journeys, nil
}

// GetDetail retrieves journey details including chapters and missions
// Returns apperror.JourneyNotFound if journey doesn't exist
func (s *JourneyService) GetDetail(ctx context.Context, journeyID int64) (*model.JourneyDetail, error) {
	// Get journey by ID
	journey, err := s.journeyRepository.GetByID(ctx, journeyID)
	if err != nil {
		if errors.Is(err, repository.ErrJourneyNotFound) {
			return nil, apperror.JourneyNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Get chapters for this journey
	chapters, err := s.journeyRepository.ListChaptersByJourneyID(ctx, journeyID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// If no chapters, return empty result
	if len(chapters) == 0 {
		return &model.JourneyDetail{
			Journey:           journey,
			Chapters:          chapters,
			MissionsByChapter: make(map[int64][]*model.Mission),
		}, nil
	}

	// Collect all chapter IDs
	chapterIDs := make([]int64, 0, len(chapters))
	for _, chapter := range chapters {
		chapterIDs = append(chapterIDs, chapter.ID)
	}

	// Get all missions for these chapters
	missions, err := s.journeyRepository.ListMissionsByChapterIDs(ctx, chapterIDs)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// Group missions by chapter_id
	missionsByChapter := make(map[int64][]*model.Mission, len(chapterIDs))
	for _, mission := range missions {
		missionsByChapter[mission.ChapterID] = append(missionsByChapter[mission.ChapterID], mission)
	}

	return &model.JourneyDetail{
		Journey:           journey,
		Chapters:          chapters,
		MissionsByChapter: missionsByChapter,
	}, nil
}
