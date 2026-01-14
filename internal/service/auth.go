package service

import (
	"context"
	"time"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
	Login(ctx context.Context, req dto.LoginRequest) (*LoginResult, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	userRepository repository.UserRepository
	tokenService   *auth.TokenService
}

// LoginResult holds the complete result of a successful login
type LoginResult struct {
	Token    string
	Expire   time.Time
	UserInfo dto.UserInfo
}

func NewAuthService(userRepository repository.UserRepository, tokenService *auth.TokenService) AuthService {
	return &authService{
		userRepository: userRepository,
		tokenService:   tokenService,
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

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*LoginResult, error) {
	// Get user by username
	user, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.AuthFailed()
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.AuthFailed()
	}

	// Generate token
	token, expire, err := s.tokenService.Generate(user)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	// Construct result
	result := &LoginResult{
		Token:  token,
		Expire: expire,
		UserInfo: dto.UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Experience: user.ExperiencePoints,
		},
	}

	return result, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	// Currently just validates token exists
	// Future: can add token blacklist logic here
	if token == "" {
		return apperror.Unauthorized()
	}
	return nil
}
