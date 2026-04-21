// Package services contains the business logic for Kodia Framework.
// Services depend only on port interfaces — never on concrete implementations.
package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/hash"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	jwtManager       *jwt.Manager
	cache            ports.CacheProvider
	mailer           ports.Mailer
	baseURL          string
	frontendURL      string
	log              *zap.Logger
}

// NewAuthService creates a new AuthService with its dependencies injected.
func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	cache ports.CacheProvider,
	mailer ports.Mailer,
	baseURL string,
	frontendURL string,
	log *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		cache:            cache,
		mailer:           mailer,
		baseURL:          baseURL,
		frontendURL:      frontendURL,
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
	role := domain.RoleUser
	_, total, err := s.userRepo.FindAll(ctx, nil)
	if err == nil && total == 0 {
		role = domain.RoleAdmin
	}

	user := &domain.User{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      role,
		IsActive:  true,
		IsVerified: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("register: %w", err)
	}

	// 4. Trigger verification email (Batteries Included)
	go func() {
		if err := s.SendVerificationEmail(context.Background(), user.ID); err != nil {
			s.log.Warn("Failed to send verification email on register", zap.Error(err))
		}
	}()

	// 5. Generate tokens
	return s.generateTokenPair(ctx, user)
}

// Login authenticates a user and returns JWT tokens or MFA requirement.
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

	// 4. Handle 2FA (Batteries Included)
	if user.TwoFactorEnabled {
		mfaToken := uuid.NewString()
		key := fmt.Sprintf("mfa_login:%s", mfaToken)
		// Store userID in cache for 5 minutes
		if err := s.cache.Set(ctx, key, user.ID, 5*time.Minute); err != nil {
			return nil, err
		}

		return &ports.AuthResponse{
			MFARequired: true,
			MFAToken:    mfaToken,
		}, nil
	}

	// 5. Self-Healing Admin Check
	if user.Role == domain.RoleUser {
		totalAdmins, err := s.userRepo.CountByRole(ctx, string(domain.RoleAdmin))
		if err == nil && totalAdmins == 0 {
			user.Role = domain.RoleAdmin
			_ = s.userRepo.Update(ctx, user)
			s.log.Info("Self-healing: Elevated first login to admin role", zap.String("email", user.Email))
		}
	}

	// 6. Generate tokens
	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) LoginVerify2FA(ctx context.Context, mfaToken string, code string) (*ports.AuthResponse, error) {
	key := fmt.Sprintf("mfa_login:%s", mfaToken)
	var userID string
	if err := s.cache.Get(ctx, key, &userID); err != nil {
		return nil, fmt.Errorf("MFA session expired or invalid")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !totp.Validate(code, user.TwoFactorSecret) {
		return nil, fmt.Errorf("invalid 2FA code")
	}

	s.cache.Delete(ctx, key)
	return s.generateTokenPair(ctx, user)
}

// RefreshToken validates a refresh token and issues a new access token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*ports.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	storedToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}
	if !storedToken.IsValid() {
		return nil, domain.ErrTokenRevoked
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !user.IsActive {
		return nil, domain.ErrInactiveAccount
	}

	if err := s.refreshTokenRepo.RevokeByToken(ctx, refreshTokenStr); err != nil {
		s.log.Warn("Failed to revoke old refresh token", zap.Error(err))
	}

	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepo.RevokeByToken(ctx, refreshToken); err != nil {
		s.log.Warn("Failed to revoke refresh token on logout", zap.Error(err))
		return nil
	}
	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
}

/* --- New Security Flows --- */

func (s *AuthService) SendVerificationEmail(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	token := uuid.NewString()
	key := fmt.Sprintf("verify_email:%s", token)
	if err := s.cache.Set(ctx, key, userID, 24*time.Hour); err != nil {
		return err
	}

	link := fmt.Sprintf("%s/auth/verify-email?token=%s", s.frontendURL, token)
	return s.mailer.Send(ctx, []string{user.Email}, "Verify Your Email - Kodia", 
		fmt.Sprintf("Welcome to the colony! Verify your email here: %s", link))
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	key := fmt.Sprintf("verify_email:%s", token)
	var userID string
	if err := s.cache.Get(ctx, key, &userID); err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.cache.Delete(ctx, key)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	token := uuid.NewString()
	key := fmt.Sprintf("password_reset:%s", token)
	if err := s.cache.Set(ctx, key, user.ID, 30*time.Minute); err != nil {
		return err
	}

	link := fmt.Sprintf("%s/auth/reset-password?token=%s", s.frontendURL, token)
	return s.mailer.Send(ctx, []string{user.Email}, "Reset Your Password - Kodia", 
		fmt.Sprintf("Forgot your sting? Reset it here: %s", link))
}

func (s *AuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	key := fmt.Sprintf("password_reset:%s", token)
	var userID string
	if err := s.cache.Get(ctx, key, &userID); err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	hashedPassword, err := hash.Make(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.cache.Delete(ctx, key)
}

func (s *AuthService) Enable2FA(ctx context.Context, userID string) (*ports.TwoFactorSetup, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Kodia Framework",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("2fa_setup:%s", userID)
	if err := s.cache.Set(ctx, cacheKey, key.Secret(), 15*time.Minute); err != nil {
		return nil, err
	}

	// Generate QR Code Image (Base64)
	qrCode, err := qr.Encode(key.URL(), qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, qrCode); err != nil {
		return nil, err
	}

	qrBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return &ports.TwoFactorSetup{
		Secret: key.Secret(),
		QRCode: qrBase64,
	}, nil
}

func (s *AuthService) Verify2FA(ctx context.Context, userID string, code string) ([]string, error) {
	cacheKey := fmt.Sprintf("2fa_setup:%s", userID)
	var secret string
	if err := s.cache.Get(ctx, cacheKey, &secret); err != nil {
		return nil, fmt.Errorf("setup session expired")
	}

	if !totp.Validate(code, secret) {
		return nil, fmt.Errorf("invalid code")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	recovery := make([]string, 8)
	for i := 0; i < 8; i++ {
		recovery[i] = uuid.NewString()[0:8]
	}

	user.TwoFactorEnabled = true
	user.TwoFactorSecret = secret
	user.TwoFactorRecoveryCodes = recovery

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	s.cache.Delete(ctx, cacheKey)
	return recovery, nil
}

func (s *AuthService) Disable2FA(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""
	user.TwoFactorRecoveryCodes = nil

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) generateTokenPair(ctx context.Context, user *domain.User) (*ports.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role), user.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenStr, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role), user.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

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
	}

	return &ports.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}
