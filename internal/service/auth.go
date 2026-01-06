package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameExists = errors.New("使用者名稱已存在")

type AuthService struct {
	queries *db.Queries
	logger  *slog.Logger
}

func NewAuthService(queries *db.Queries, logger *slog.Logger) *AuthService {
	return &AuthService{
		queries: queries,
		logger:  logger,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// Check if username already exists
	_, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err == nil {
		// User found, username exists
		return 0, ErrUsernameExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		// Database error
		s.logger.Error("Failed to check username existence", "error", err, "username", req.Username)
		return 0, err
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return 0, err
	}

	// Create user
	userID, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		s.logger.Error("Failed to create user", "error", err, "username", req.Username)
		return 0, err
	}

	s.logger.Info("User registered successfully", "userId", userID, "username", req.Username)
	return userID, nil
}
