package auth

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// Token type constants
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// TokenGenerator defines the interface for JWT token generation and verification
type TokenGenerator interface {
	// GenerateAccessToken creates a new access token for the user
	GenerateAccessToken(user *model.User) (token string, jti string, expiresAt time.Time, err error)
	// GenerateRefreshToken creates a new refresh token for the user
	GenerateRefreshToken(userID int64) (token string, jti string, expiresAt time.Time, err error)
	// ParseAccessToken extracts JTI, user ID, and expiry time from an access token
	ParseAccessToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error)
	// ParseRefreshToken extracts JTI, user ID, and expiry time from a refresh token
	ParseRefreshToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error)
}

// jwtTokenGenerator implements TokenGenerator interface
type jwtTokenGenerator struct {
	config config.JWTConfig
}

// NewTokenGenerator creates a new TokenGenerator instance
func NewTokenGenerator(config config.JWTConfig) TokenGenerator {
	return &jwtTokenGenerator{config: config}
}

// GenerateAccessToken creates a new access token for the user
func (tg *jwtTokenGenerator) GenerateAccessToken(user *model.User) (string, string, time.Time, error) {
	expireTime := time.Now().Add(tg.config.AccessTokenTimeout)
	jti := uuid.New().String()

	claims := gojwt.MapClaims{
		"jti":      jti,
		"type":     TokenTypeAccess,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expireTime.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tg.config.Secret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, jti, expireTime, nil
}

// GenerateRefreshToken creates a new refresh token for the user
func (tg *jwtTokenGenerator) GenerateRefreshToken(userID int64) (string, string, time.Time, error) {
	expireTime := time.Now().Add(tg.config.RefreshTokenTimeout)
	jti := uuid.New().String()

	claims := gojwt.MapClaims{
		"jti":     jti,
		"type":    TokenTypeRefresh,
		"user_id": userID,
		"exp":     expireTime.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tg.config.Secret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, jti, expireTime, nil
}

// ParseAccessToken extracts JTI, user ID, and expiry time from an access token
func (tg *jwtTokenGenerator) ParseAccessToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error) {
	return tg.parseToken(tokenString, TokenTypeAccess)
}

// ParseRefreshToken extracts JTI, user ID, and expiry time from a refresh token
func (tg *jwtTokenGenerator) ParseRefreshToken(tokenString string) (jti string, userID int64, expiresAt time.Time, err error) {
	return tg.parseToken(tokenString, TokenTypeRefresh)
}

// parseToken is a helper function to parse and validate a token with type checking
func (tg *jwtTokenGenerator) parseToken(tokenString string, expectedType string) (jti string, userID int64, expiresAt time.Time, err error) {
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

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok {
		// For backward compatibility, treat tokens without type as access tokens
		if expectedType == TokenTypeRefresh {
			return "", 0, time.Time{}, fmt.Errorf("invalid token type: expected %s", expectedType)
		}
	} else if tokenType != expectedType {
		return "", 0, time.Time{}, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, tokenType)
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
