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

	// CountByRole returns the total number of users with a specific role.
	CountByRole(ctx context.Context, role string) (int64, error)
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

// NotificationRepository defines all database operations for the Notification entity.
type NotificationRepository interface {
	// Create persists a new notification to the database.
	Create(ctx context.Context, n *domain.Notification) error

	// FindByID retrieves a notification by its unique ID.
	FindByID(ctx context.Context, id string) (*domain.Notification, error)

	// FindByUserID retrieves paginated notifications for a specific user.
	FindByUserID(ctx context.Context, userID string, params *pagination.Params) ([]*domain.Notification, int64, error)

	// MarkAsRead marks a single notification as read.
	MarkAsRead(ctx context.Context, id string, userID string) error

	// MarkAllAsRead marks all notifications for a user as read.
	MarkAllAsRead(ctx context.Context, userID string) error

	// Delete removes a notification (ownership checked via userID).
	Delete(ctx context.Context, id string, userID string) error

	// CountUnread returns the number of unread notifications for a user.
	CountUnread(ctx context.Context, userID string) (int64, error)
}

// RoleRepository defines all database operations for roles.
type RoleRepository interface {
	// Create persists a new role to the database.
	Create(ctx context.Context, role *domain.RoleEntity) error

	// FindByName retrieves a role by its name.
	FindByName(ctx context.Context, name string) (*domain.RoleEntity, error)

	// FindAll retrieves all roles.
	FindAll(ctx context.Context) ([]*domain.RoleEntity, error)

	// Update persists changes to an existing role.
	Update(ctx context.Context, role *domain.RoleEntity) error

	// Delete soft-deletes a role by its ID.
	Delete(ctx context.Context, id string) error

	// AssignToUser assigns a role to a user.
	AssignToUser(ctx context.Context, userID, roleName string) error

	// RevokeFromUser revokes a role from a user.
	RevokeFromUser(ctx context.Context, userID, roleName string) error

	// GetUserRoles retrieves all role names assigned to a user.
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
}

// PermissionRepository defines all database operations for permissions.
type PermissionRepository interface {
	// FindAll retrieves all permissions.
	FindAll(ctx context.Context) ([]*domain.PermissionEntity, error)

	// FindByName retrieves a permission by its name.
	FindByName(ctx context.Context, name string) (*domain.PermissionEntity, error)

	// Create persists a new permission to the database.
	Create(ctx context.Context, perm *domain.PermissionEntity) error

	// FindByGroup retrieves all permissions in a group.
	FindByGroup(ctx context.Context, group string) ([]*domain.PermissionEntity, error)
}

