package authsocial

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/hash"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"go.uber.org/zap"
)

// AuthResult holds the result of successful social authentication.
type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
}

// SocialAuthService handles the OAuth callback flow.
type SocialAuthService struct {
	socialRepo SocialAccountRepository
	userRepo   ports.UserRepository
	cache      ports.CacheProvider
	jwtManager *jwt.Manager
	baseURL    string
	log        *zap.Logger
}

// NewSocialAuthService creates a new social auth service.
func NewSocialAuthService(
	socialRepo SocialAccountRepository,
	userRepo ports.UserRepository,
	cache ports.CacheProvider,
	jwtManager *jwt.Manager,
	baseURL string,
	log *zap.Logger,
) *SocialAuthService {
	return &SocialAuthService{
		socialRepo: socialRepo,
		userRepo:   userRepo,
		cache:      cache,
		jwtManager: jwtManager,
		baseURL:    baseURL,
		log:        log,
	}
}

// GenerateState creates a CSRF protection state and stores it in cache.
func (s *SocialAuthService) GenerateState(ctx context.Context) (string, error) {
	state := uuid.New().String()
	key := fmt.Sprintf("social_state:%s", state)
	if err := s.cache.Set(ctx, key, "1", 5*time.Minute); err != nil {
		s.log.Error("failed to store state in cache", zap.Error(err))
		return "", err
	}
	return state, nil
}

// VerifyState checks if the state is valid and removes it from cache.
func (s *SocialAuthService) VerifyState(ctx context.Context, state string) bool {
	key := fmt.Sprintf("social_state:%s", state)
	val, err := s.cache.Get(ctx, key)
	if err != nil || val == "" {
		return false
	}
	_ = s.cache.Delete(ctx, key) // Clean up
	return true
}

// HandleCallback handles the OAuth callback and returns auth tokens.
func (s *SocialAuthService) HandleCallback(
	ctx context.Context,
	provider string,
	code string,
	state string,
	oauthProvider Provider,
) (*AuthResult, error) {
	// 1. Verify CSRF state
	if !s.VerifyState(ctx, state) {
		return nil, errors.New("invalid or expired state parameter")
	}

	// 2. Exchange code for token
	token, err := oauthProvider.Exchange(ctx, code)
	if err != nil {
		s.log.Error("failed to exchange code", zap.Error(err))
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	// 3. Get user info from provider
	socialUser, err := oauthProvider.GetUser(ctx, token)
	if err != nil {
		s.log.Error("failed to get user from provider", zap.Error(err))
		return nil, fmt.Errorf("failed to get user information: %w", err)
	}

	// 4. Check if social account exists
	existingAccount, err := s.socialRepo.FindByProvider(ctx, provider, socialUser.ID)
	if err != nil {
		s.log.Error("failed to query social account", zap.Error(err))
		return nil, err
	}

	var userID string

	if existingAccount != nil {
		// Account exists, use existing user
		userID = existingAccount.UserID
		s.log.Info("existing social account found", zap.String("user_id", userID))
	} else {
		// 5. Check if user exists by email
		existingUser, err := s.userRepo.FindByEmail(ctx, socialUser.Email)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			s.log.Error("failed to query user by email", zap.Error(err))
			return nil, err
		}

		if existingUser != nil {
			// User exists, link social account
			userID = existingUser.ID
			s.log.Info("user found by email, linking social account", zap.String("user_id", userID))
		} else {
			// 6. Create new user
			userID = uuid.New().String()
			newUser := &domain.User{
				ID:         userID,
				Name:       socialUser.Name,
				Email:      socialUser.Email,
				Password:   "", // No password for social logins
				Role:       domain.RoleUser,
				IsActive:   true,
				IsVerified: true, // Trust social providers for verification
				AvatarURL:  &socialUser.AvatarURL,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := s.userRepo.Create(ctx, newUser); err != nil {
				s.log.Error("failed to create user", zap.Error(err))
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			s.log.Info("new user created from social login", zap.String("user_id", userID))
		}

		// 7. Create social account link
		socialAccount := &SocialAccount{
			ID:         uuid.New().String(),
			UserID:     userID,
			Provider:   provider,
			ProviderID: socialUser.ID,
			Email:      socialUser.Email,
			Name:       socialUser.Name,
			AvatarURL:  socialUser.AvatarURL,
			CreatedAt:  time.Now(),
		}

		if err := s.socialRepo.Create(ctx, socialAccount); err != nil {
			s.log.Error("failed to create social account", zap.Error(err))
			return nil, fmt.Errorf("failed to link social account: %w", err)
		}

		s.log.Info("social account created", zap.String("user_id", userID))
	}

	// 8. Get user from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to retrieve user", zap.Error(err))
		return nil, err
	}

	// 9. Generate JWT tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		s.log.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		s.log.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
