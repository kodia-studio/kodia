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
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	TokenType   TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenStore defines the interface for tracking token lifecycle (revocation and reuse).
type TokenStore interface {
	// IsRevoked checks if a token ID has been explicitly revoked.
	IsRevoked(jti string) (bool, error)
	// MarkUsed marks a refresh token as used and detects if it was already used (Reuse Detection).
	// Returns (alreadyUsed, error).
	MarkUsed(jti string, userID string, expiry time.Duration) (bool, error)
	// Revoke sets a token (or all tokens for a user) as invalid.
	Revoke(jti string, expiry time.Duration) error
}

// Manager manages JWT creation and verification.
type Manager struct {
	accessSecret       []byte
	refreshSecret      []byte
	accessExpiryHours  time.Duration
	refreshExpiryHours time.Duration
	store              TokenStore
}

// NewManager creates a new JWT Manager.
func NewManager(accessSecret, refreshSecret string, accessExpiryHours, refreshExpiryDays int) *Manager {
	return &Manager{
		accessSecret:       []byte(accessSecret),
		refreshSecret:      []byte(refreshSecret),
		accessExpiryHours:  time.Duration(accessExpiryHours) * time.Hour,
		refreshExpiryHours: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

// SetStore sets the token store for revocation and rotation support.
func (m *Manager) SetStore(store TokenStore) {
	m.store = store
}

// GenerateAccessToken creates a short-lived access token.
func (m *Manager) GenerateAccessToken(userID, email, role string, permissions []string) (string, error) {
	return m.generate(userID, email, role, permissions, AccessToken, m.accessSecret, m.accessExpiryHours)
}

// GenerateRefreshToken creates a long-lived refresh token.
func (m *Manager) GenerateRefreshToken(userID, email, role string, permissions []string) (string, error) {
	return m.generate(userID, email, role, permissions, RefreshToken, m.refreshSecret, m.refreshExpiryHours)
}

func (m *Manager) generate(userID, email, role string, permissions []string, tokenType TokenType, secret []byte, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:      userID,
		Email:       email,
		Role:        role,
		Permissions: permissions,
		TokenType:   tokenType,
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
// This method also handles rotation and reuse detection if a Store is configured.
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.validate(tokenString, m.refreshSecret, RefreshToken)
	if err != nil {
		return nil, err
	}

	// If a store is configured, check for revocation and reuse
	if m.store != nil {
		// 1. Check if revoked
		revoked, _ := m.store.IsRevoked(claims.ID)
		if revoked {
			return nil, errors.New("token has been revoked")
		}

		// 2. Mark as used and detect reuse
		// The caller should ideally handle the "rotate" part, but we provide the detection here.
		alreadyUsed, err := m.store.MarkUsed(claims.ID, claims.UserID, claims.ExpiresAt.Time.Sub(time.Now()))
		if err != nil {
			return nil, err
		}
		if alreadyUsed {
			// REUSE DETECTED: This is a critical security event.
			// Standard behavior: Invalidate all of the user's refresh tokens.
			_ = m.store.Revoke(claims.UserID, m.refreshExpiryHours) // Assuming store supports revoking by UserID
			return nil, errors.New("refresh token reuse detected - all sessions invalidated")
		}
	}

	return claims, nil
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
