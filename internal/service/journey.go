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

// journeyProgressRepository provides user mission progress data access operations for JourneyService
type journeyProgressRepository interface {
	ListByUserAndMissions(ctx context.Context, userID int64, missionIDs []int64) (map[int64]string, error)
}

// JourneyService implements journey business logic operations.
type JourneyService struct {
	journeyRepository  journeyRepository
	progressRepository journeyProgressRepository
}

// NewJourneyService creates a new JourneyService instance.
func NewJourneyService(journeyRepository *repository.JourneyRepository, progressRepository *repository.ProgressRepository) *JourneyService {
	return &JourneyService{
		journeyRepository:  journeyRepository,
		progressRepository: progressRepository,
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
// If userID is provided (authenticated user), also fetches mission progress
// Returns apperror.JourneyNotFound if journey doesn't exist
func (s *JourneyService) GetDetail(ctx context.Context, journeyID int64, userID *int64) (*model.JourneyDetail, error) {
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
			ProgressByMission: nil,
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

	// Get user progress if authenticated
	var progressByMission map[int64]string
	if userID != nil && len(missions) > 0 {
		// Collect all mission IDs
		missionIDs := make([]int64, 0, len(missions))
		for _, mission := range missions {
			missionIDs = append(missionIDs, mission.ID)
		}

		// Fetch progress for all missions
		progressByMission, err = s.progressRepository.ListByUserAndMissions(ctx, *userID, missionIDs)
		if err != nil {
			return nil, apperror.DatabaseError(err)
		}
	}

	return &model.JourneyDetail{
		Journey:           journey,
		Chapters:          chapters,
		MissionsByChapter: missionsByChapter,
		ProgressByMission: progressByMission,
	}, nil
}
