// Package ports defines the interfaces (contracts) between layers in Kodia Framework.
// The core domain depends only on these interfaces — never on concrete implementations.
// This enforces the Dependency Inversion Principle and makes testing trivial.
package ports

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/pagination"
)

// --- Repository Interfaces ---

// UserRepository defines all database operations for the User entity.
// Implemented in internal/adapters/repository/postgres/ or /mysql/.
type UserRepository interface {
	// Create persists a new user to the database.
	Create(ctx context.Context, user *domain.User) error

	// FindByID retrieves a user by their unique ID.
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// FindAll retrieves a paginated list of users.
	FindAll(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error)

	// Update persists changes to an existing user.
	Update(ctx context.Context, user *domain.User) error

	// Delete soft-deletes a user by their ID.
	Delete(ctx context.Context, id string) error

	// ExistsByEmail returns true if a user with the given email exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// RefreshTokenRepository defines all database operations for refresh tokens.
type RefreshTokenRepository interface {
	// Create persists a new refresh token.
	Create(ctx context.Context, token *domain.RefreshToken) error

	// FindByToken retrieves a refresh token by its token string.
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)

	// RevokeByToken marks a specific refresh token as revoked.
	RevokeByToken(ctx context.Context, token string) error

	// RevokeAllForUser revokes all active refresh tokens for a user (logout from all devices).
	RevokeAllForUser(ctx context.Context, userID string) error

	// DeleteExpired removes all expired refresh tokens (for cleanup).
	DeleteExpired(ctx context.Context) error
}

// ProductRepository defines all database operations for the Product entity.
type ProductRepository interface {
	// Create persists a new product to the database.
	Create(ctx context.Context, product *domain.Product) error

	// FindByID retrieves a product by its unique ID.
	FindByID(ctx context.Context, id string) (*domain.Product, error)

	// FindAll retrieves a paginated list of products.
	FindAll(ctx context.Context, params *pagination.Params) ([]domain.Product, int64, error)

	// Update persists changes to an existing product.
	Update(ctx context.Context, product *domain.Product) error

	// Delete soft-deletes a product by its ID.
	Delete(ctx context.Context, id string) error
}

