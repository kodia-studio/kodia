package ports

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/pagination"
)

// --- Service Interfaces ---

// AuthService defines the authentication business operations.
// Implemented in internal/core/services/auth_service.go.
type AuthService interface {
	// Register creates a new user account.
	Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)

	// Login authenticates a user and returns JWT tokens or MFA requirement.
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)

	// RefreshToken generates a new access token from a valid refresh token.
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)

	// Logout revokes the provided refresh token.
	Logout(ctx context.Context, refreshToken string) error

	// LogoutAll revokes all refresh tokens for a user.
	LogoutAll(ctx context.Context, userID string) error

	// Email Verification
	VerifyEmail(ctx context.Context, token string) error
	SendVerificationEmail(ctx context.Context, userID string) error

	// Password Recovery
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error

	// 2FA Security
	Enable2FA(ctx context.Context, userID string) (*TwoFactorSetup, error)
	Verify2FA(ctx context.Context, userID string, code string) ([]string, error)
	Disable2FA(ctx context.Context, userID string) error
	LoginVerify2FA(ctx context.Context, mfaToken string, code string) (*AuthResponse, error)
}

type TwoFactorSetup struct {
	Secret string
	QRCode string // Base64 or URL
}

// UserService defines user management business operations.
// Implemented in internal/core/services/user_service.go.
type UserService interface {
	// GetByID fetches a single user by their ID.
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetAll returns a paginated list of users.
	GetAll(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error)

	// Update updates a user's profile information.
	Update(ctx context.Context, id string, input UpdateUserInput) (*domain.User, error)

	// Delete soft-deletes a user.
	Delete(ctx context.Context, id string) error

	// ChangePassword updates a user's password after verifying the current one.
	ChangePassword(ctx context.Context, id string, input ChangePasswordInput) error

	// UpdateAvatar updates the avatar URL for a user.
	UpdateAvatar(ctx context.Context, id string, avatarURL string) error
}

// ProductService defines product management business operations.
type ProductService interface {
	// GetAll returns a paginated list of products.
	GetAll(ctx context.Context, params *pagination.Params) ([]domain.Product, int64, error)

	// GetByID fetches a single product by its ID.
	GetByID(ctx context.Context, id string) (*domain.Product, error)

	// Delete soft-deletes a product.
	Delete(ctx context.Context, id string) error
}


// --- Input/Output DTOs for Services ---
// These are separate from HTTP DTOs — they represent the contract between HTTP and Service layers.

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
	MFARequired  bool   `json:"mfa_required,omitempty"`
	MFAToken     string `json:"mfa_token,omitempty"`
}

type UpdateUserInput struct {
	Name      *string
	AvatarURL *string
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

// NotificationService defines notification business operations.
type NotificationService interface {
	// Send creates, persists, and pushes a notification to the user.
	Send(ctx context.Context, input SendNotificationInput) (*domain.Notification, error)

	// GetAll returns paginated notifications for a user.
	GetAll(ctx context.Context, userID string, params *pagination.Params) ([]*domain.Notification, int64, error)

	// MarkAsRead marks a single notification as read.
	MarkAsRead(ctx context.Context, id string, userID string) error

	// MarkAllAsRead marks all notifications for a user as read.
	MarkAllAsRead(ctx context.Context, userID string) error

	// Delete removes a notification.
	Delete(ctx context.Context, id string, userID string) error

	// CountUnread returns the number of unread notifications.
	CountUnread(ctx context.Context, userID string) (int64, error)
}

// SendNotificationInput is the input for sending a notification.
type SendNotificationInput struct {
	UserID    string
	Type      domain.NotificationType
	Title     string
	Message   string
	Data      map[string]interface{}
	SendEmail bool
}
