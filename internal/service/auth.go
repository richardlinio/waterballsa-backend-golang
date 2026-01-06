package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameExists = errors.New("使用者名稱已存在")

type AuthService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to check username existence", "error", err, "username", req.Username)
		return 0, err
	}
	if exists {
		return 0, ErrUsernameExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return 0, err
	}

	// Create user
	userID, err := s.userRepo.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		s.logger.Error("Failed to create user", "error", err, "username", req.Username)
		return 0, err
	}

	return userID, nil
}
