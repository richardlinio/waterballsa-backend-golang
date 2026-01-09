package service

import (
	"context"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// bcrypt can only handle up to 72 bytes
	if len(req.Password) > 72 {
		return 0, apperror.PasswordTooLong()
	}

	// Check if username already exists
	exists, err := s.userRepository.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return 0, apperror.DatabaseError(err)
	}
	if exists {
		return 0, apperror.UsernameExists()
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, apperror.InternalError(err)
	}

	// Create user
	userID, err := s.userRepository.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		return 0, apperror.DatabaseError(err)
	}

	return userID, nil
}
