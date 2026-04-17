package ports

import (
	"context"

	"github.com/kodia/framework/backend/internal/core/domain"
	"github.com/kodia/framework/backend/pkg/pagination"
)

// --- Service Interfaces ---

// AuthService defines the authentication business operations.
// Implemented in internal/core/services/auth_service.go.
type AuthService interface {
	// Register creates a new user account.
	Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)

	// Login authenticates a user and returns JWT tokens.
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)

	// RefreshToken generates a new access token from a valid refresh token.
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)

	// Logout revokes the provided refresh token.
	Logout(ctx context.Context, refreshToken string) error

	// LogoutAll revokes all refresh tokens for a user.
	LogoutAll(ctx context.Context, userID string) error
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
}

type UpdateUserInput struct {
	Name      *string
	AvatarURL *string
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}
