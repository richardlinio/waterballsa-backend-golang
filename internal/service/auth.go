package service

import (
	"context"
	"errors"
	"time"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// userRepository provides user data access operations for AuthService
type userRepository interface {
	Create(ctx context.Context, username, passwordHash string) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// accessTokenRepository provides access token blacklist operations for AuthService
type accessTokenRepository interface {
	Invalidate(ctx context.Context, jti string, userID int64, expiresAt time.Time) error
	IsInvalidated(ctx context.Context, jti string) (bool, error)
}

// refreshTokenRepository provides refresh token persistence operations for AuthService
type refreshTokenRepository interface {
	Create(ctx context.Context, jti string, userID int64, expiresAt time.Time) error
	GetByJTI(ctx context.Context, jti string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, jti string) error
}

// tokenGenerator provides JWT token generation and parsing operations for AuthService
type tokenGenerator interface {
	GenerateAccessToken(user *model.User) (token string, jti string, expiresAt time.Time, err error)
	GenerateRefreshToken(userID int64) (token string, jti string, expiresAt time.Time, err error)
	ParseAccessToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error)
	ParseRefreshToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error)
}

// AuthService implements authentication business logic operations.
type AuthService struct {
	userRepository         userRepository
	accessTokenRepository  accessTokenRepository
	refreshTokenRepository refreshTokenRepository
	tokenGenerator         tokenGenerator
}

// LoginResult holds the complete result of a successful login
type LoginResult struct {
	AccessToken        string
	AccessTokenExpire  time.Time
	RefreshToken       string
	RefreshTokenExpire time.Time
	UserInfo           dto.UserInfo
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(
	userRepository *repository.UserRepository,
	accessTokenRepository *repository.AccessTokenRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	tokenGenerator *auth.JWTTokenGenerator,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		tokenGenerator:         tokenGenerator,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
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
		return 0, apperror.InternalServerError(err)
	}

	// Create user
	userID, err := s.userRepository.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		return 0, apperror.DatabaseError(err)
	}

	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*LoginResult, error) {
	// Get user by username
	user, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.AuthFailed()
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.AuthFailed()
	}

	// Generate access token
	accessToken, _, accessExpire, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// Generate refresh token
	refreshToken, refreshJTI, refreshExpire, err := s.tokenGenerator.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// Store refresh token in database
	if err := s.refreshTokenRepository.Create(ctx, refreshJTI, user.ID, refreshExpire); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &LoginResult{
		AccessToken:        accessToken,
		AccessTokenExpire:  accessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshExpire,
		UserInfo: dto.UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Experience: user.ExperiencePoints,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// Invalidate access token if provided
	if accessToken != "" {
		jti, userID, expiresAt, err := s.tokenGenerator.ParseAccessToken(accessToken)
		if err == nil {
			// Best effort: add to blacklist, ignore errors
			_ = s.accessTokenRepository.Invalidate(ctx, jti, userID, expiresAt)
		}
	}

	// Revoke refresh token if provided
	if refreshToken != "" {
		jti, _, _, err := s.tokenGenerator.ParseRefreshToken(refreshToken)
		if err == nil {
			// Best effort: revoke refresh token, ignore errors
			_ = s.refreshTokenRepository.Revoke(ctx, jti)
		}
	}

	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	if refreshToken == "" {
		return nil, apperror.Unauthorized()
	}

	// Parse refresh token
	jti, userID, _, err := s.tokenGenerator.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.Unauthorized()
	}

	// Check if refresh token exists and is valid in database
	dbToken, err := s.refreshTokenRepository.GetByJTI(ctx, jti)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, apperror.Unauthorized()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Verify user ID matches
	if dbToken.UserID != userID {
		return nil, apperror.Unauthorized()
	}

	// Get fresh user data
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Unauthorized()
	}

	// Revoke the old refresh token (token rotation)
	if err := s.refreshTokenRepository.Revoke(ctx, jti); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// Generate new access token
	newAccessToken, _, accessExpire, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// Generate new refresh token
	newRefreshToken, newRefreshJTI, refreshExpire, err := s.tokenGenerator.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// Store new refresh token in database
	if err := s.refreshTokenRepository.Create(ctx, newRefreshJTI, user.ID, refreshExpire); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &LoginResult{
		AccessToken:        newAccessToken,
		AccessTokenExpire:  accessExpire,
		RefreshToken:       newRefreshToken,
		RefreshTokenExpire: refreshExpire,
		UserInfo: dto.UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Experience: user.ExperiencePoints,
		},
	}, nil
}
