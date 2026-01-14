package auth

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// TokenGenerator defines the interface for JWT token generation and verification
type TokenGenerator interface {
	Generate(user *model.User) (string, time.Time, error)
}

// jwtTokenGenerator implements TokenGenerator interface
type jwtTokenGenerator struct {
	config config.JWTConfig
}

// NewTokenGenerator creates a new TokenGenerator instance
func NewTokenGenerator(config config.JWTConfig) TokenGenerator {
	return &jwtTokenGenerator{config: config}
}

// Generate creates a new JWT token for the user
func (tg *jwtTokenGenerator) Generate(user *model.User) (string, time.Time, error) {
	expireTime := time.Now().Add(tg.config.AccessTokenTimeout)

	claims := gojwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"experience": user.ExperiencePoints,
		"exp":        expireTime.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tg.config.Secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expireTime, nil
}
