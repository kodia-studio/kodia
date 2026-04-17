// Package jwt provides JWT access and refresh token utilities for Kodia Framework.
// Uses golang-jwt/jwt v5 with HS256 signing by default.
// Designed to be extensible for RS256/ES256 in production environments.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType distinguishes access tokens from refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT payload for Kodia Framework.
type Claims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Manager manages JWT creation and verification.
type Manager struct {
	accessSecret       []byte
	refreshSecret      []byte
	accessExpiryHours  time.Duration
	refreshExpiryHours time.Duration
}

// NewManager creates a new JWT Manager.
// accessSecret and refreshSecret should be different strong secrets (32+ chars).
func NewManager(accessSecret, refreshSecret string, accessExpiryHours, refreshExpiryDays int) *Manager {
	return &Manager{
		accessSecret:       []byte(accessSecret),
		refreshSecret:      []byte(refreshSecret),
		accessExpiryHours:  time.Duration(accessExpiryHours) * time.Hour,
		refreshExpiryHours: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

// GenerateAccessToken creates a short-lived access token.
func (m *Manager) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generate(userID, email, role, AccessToken, m.accessSecret, m.accessExpiryHours)
}

// GenerateRefreshToken creates a long-lived refresh token.
func (m *Manager) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generate(userID, email, role, RefreshToken, m.refreshSecret, m.refreshExpiryHours)
}

func (m *Manager) generate(userID, email, role string, tokenType TokenType, secret []byte, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "kodia",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateAccessToken verifies an access token and returns its claims.
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validate(tokenString, m.accessSecret, AccessToken)
}

// ValidateRefreshToken verifies a refresh token and returns its claims.
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validate(tokenString, m.refreshSecret, RefreshToken)
}

func (m *Manager) validate(tokenString string, secret []byte, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
