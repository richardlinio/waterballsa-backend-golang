package auth

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// TokenGenerator defines the interface for JWT token generation and verification
type TokenGenerator interface {
	Generate(user *model.User) (string, time.Time, error)
	ParseToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error)
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
	jti := uuid.New().String()

	claims := gojwt.MapClaims{
		"jti":      jti,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expireTime.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tg.config.Secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expireTime, nil
}

// ParseToken extracts JTI, user ID, and expiry time from a token string
func (tg *jwtTokenGenerator) ParseToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tg.config.Secret, nil
	})
	if err != nil {
		return "", 0, time.Time{}, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return "", 0, time.Time{}, fmt.Errorf("invalid token claims")
	}

	jtiValue, ok := claims["jti"].(string)
	if !ok {
		return "", 0, time.Time{}, fmt.Errorf("jti claim not found or invalid")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, time.Time{}, fmt.Errorf("user_id claim not found or invalid")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", 0, time.Time{}, fmt.Errorf("exp claim not found or invalid")
	}

	return jtiValue, int64(userIDFloat), time.Unix(int64(expFloat), 0), nil
}
