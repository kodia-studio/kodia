// Package services contains the business logic for Kodia Framework.
// Services depend only on port interfaces — never on concrete implementations.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/hash"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"go.uber.org/zap"
)

// AuthService implements ports.AuthService.
type AuthService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	jwtManager       *jwt.Manager
	log              *zap.Logger
}

// NewAuthService creates a new AuthService with its dependencies injected.
func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	log *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		log:              log,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, input ports.RegisterInput) (*ports.AuthResponse, error) {
	// 1. Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		s.log.Error("Failed to check email existence", zap.Error(err))
		return nil, fmt.Errorf("register: %w", err)
	}
	if exists {
		return nil, domain.ErrAlreadyExists
	}

	// 2. Hash password
	hashedPassword, err := hash.Make(input.Password)
	if err != nil {
		return nil, fmt.Errorf("register: failed to hash password: %w", err)
	}

	// 3. Create user entity
	user := &domain.User{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      domain.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("register: %w", err)
	}

	// 4. Generate tokens
	return s.generateTokenPair(ctx, user)
}

// Login authenticates a user and returns JWT tokens.
func (s *AuthService) Login(ctx context.Context, input ports.LoginInput) (*ports.AuthResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: %w", err)
	}

	// 2. Check account status
	if !user.IsActive {
		return nil, domain.ErrInactiveAccount
	}

	// 3. Verify password
	if !hash.Check(input.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// 4. Generate tokens
	return s.generateTokenPair(ctx, user)
}

// RefreshToken validates a refresh token and issues a new access token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*ports.AuthResponse, error) {
	// 1. Validate the refresh token JWT structure
	claims, err := s.jwtManager.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	// 2. Check if the token is still active in the database
	storedToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}
	if !storedToken.IsValid() {
		return nil, domain.ErrTokenRevoked
	}

	// 3. Fetch the user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !user.IsActive {
		return nil, domain.ErrInactiveAccount
	}

	// 4. Rotate token: revoke old, issue new pair
	if err := s.refreshTokenRepo.RevokeByToken(ctx, refreshTokenStr); err != nil {
		s.log.Warn("Failed to revoke old refresh token", zap.Error(err))
	}

	return s.generateTokenPair(ctx, user)
}

// Logout revokes a specific refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepo.RevokeByToken(ctx, refreshToken); err != nil {
		s.log.Warn("Failed to revoke refresh token on logout", zap.Error(err))
		return nil // Treat as success to avoid leaking info
	}
	return nil
}

// LogoutAll revokes all refresh tokens for a user.
func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
}

// generateTokenPair creates and stores an access + refresh token pair.
func (s *AuthService) generateTokenPair(ctx context.Context, user *domain.User) (*ports.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role), user.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenStr, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role), user.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Persist refresh token
	refreshToken := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		s.log.Error("Failed to persist refresh token", zap.Error(err))
		// Non-fatal: access token is still valid
	}

	return &ports.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}
