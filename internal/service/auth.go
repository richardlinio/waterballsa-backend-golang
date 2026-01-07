package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists  = errors.New("username already exists")
	ErrPasswordTooLong = errors.New("password too long")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
}

type authService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, logger *slog.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// bcrypt can only handle up to 72 bytes
	if len(req.Password) > 72 {
		return 0, ErrPasswordTooLong
	}

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to check username existence", "error", err, "username", req.Username)
		return 0, fmt.Errorf("failed to check username existence [username=%s]: %w", req.Username, err)
	}
	if exists {
		return 0, ErrUsernameExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	userID, err := s.userRepo.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		s.logger.Error("Failed to create user", "error", err, "username", req.Username)
		return 0, fmt.Errorf("failed to create user [username=%s]: %w", req.Username, err)
	}

	return userID, nil
}
