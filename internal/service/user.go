package service

import (
	"context"
	"errors"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

type userProfileRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type UserService struct {
	userProfileRepository userProfileRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userProfileRepository: userRepository,
	}
}

// GetProfile retrieves user profile by user ID
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userProfileRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperror.Unauthorized()
		}
		return nil, apperror.DatabaseError(err)
	}

	return user, nil
}
