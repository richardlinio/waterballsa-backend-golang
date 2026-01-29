package service

import (
	"context"
	"errors"
	"time"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
	"github.com/richardlinio/waterballsa-backend-golang/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// authUserRepository provides user data access operations for AuthService
type authUserRepository interface {
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
	userRepository         authUserRepository
	accessTokenRepository  accessTokenRepository
	refreshTokenRepository refreshTokenRepository
	tokenGenerator         tokenGenerator
	store                  *store.Store
	transactionTimeout     time.Duration
}

// LoginResult holds the complete result of a successful login
type LoginResult struct {
	AccessToken        string
	AccessTokenExpire  time.Time
	RefreshToken       string
	RefreshTokenExpire time.Time
	UserID             int64
	Username           string
	Experience         int32
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(
	userRepository *repository.UserRepository,
	accessTokenRepository *repository.AccessTokenRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	tokenGenerator *auth.JWTTokenGenerator,
	st *store.Store,
	transactionTimeout time.Duration,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		tokenGenerator:         tokenGenerator,
		store:                  st,
		transactionTimeout:     transactionTimeout,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// bcrypt can only handle up to 72 bytes
	if len(req.Password) > 72 {
		return 0, apperror.PasswordTooLong()
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, apperror.InternalServerError(err)
	}

	// Create user
	userID, err := s.userRepository.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		// Handle duplicate username error
		if errors.Is(err, repository.ErrUsernameDuplicate) {
			return 0, apperror.UsernameExists()
		}
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
		UserID:             user.ID,
		Username:           user.Username,
		Experience:         user.ExperiencePoints,
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

	// Parse refresh token (outside transaction - CPU-bound)
	jti, userID, _, err := s.tokenGenerator.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.Unauthorized()
	}

	// Get fresh user data (outside transaction - for response only)
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Unauthorized()
	}

	// Generate new tokens (outside transaction - CPU-bound, fail-fast)
	newAccessToken, _, accessExpire, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	newRefreshToken, newRefreshJTI, refreshExpire, err := s.tokenGenerator.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// Execute token rotation transaction
	// This ensures atomic revocation and creation of tokens
	txCtx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	_, err = s.store.RefreshTokenTx(txCtx, store.RefreshTokenTxParams{
		OldJTI:       jti,
		UserID:       userID,
		NewJTI:       newRefreshJTI,
		NewExpiresAt: refreshExpire,
	})
	if err != nil {
		// Check for specific error types - all auth failures should return Unauthorized
		if errors.Is(err, repository.ErrRefreshTokenNotFound) || errors.Is(err, store.ErrTokenUserMismatch) {
			return nil, apperror.Unauthorized()
		}
		return nil, apperror.DatabaseError(err)
	}

	return &LoginResult{
		AccessToken:        newAccessToken,
		AccessTokenExpire:  accessExpire,
		RefreshToken:       newRefreshToken,
		RefreshTokenExpire: refreshExpire,
		UserID:             user.ID,
		Username:           user.Username,
		Experience:         user.ExperiencePoints,
	}, nil
}
