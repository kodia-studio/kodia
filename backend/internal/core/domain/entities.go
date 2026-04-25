// Package domain contains the core entities and value objects of Kodia Framework.
// This package has NO external dependencies — only pure Go stdlib.
package domain

import (
	"time"
)

// UserRole represents the role of a user in the system.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User is the core user entity.
type User struct {
	ID        string
	Name      string
	Email     string
	Password  string // hashed
	Role      UserRole
	Roles     []string  // RBAC: role names assigned to user
	Permissions []string
	IsActive    bool
	IsVerified  bool

	// 2FA Security
	TwoFactorEnabled      bool
	TwoFactorSecret       string
	TwoFactorRecoveryCodes []string

	AvatarURL *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Product is the core product entity.
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// IsAdmin returns true if the user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Can checks if the user has the specified permission.
// Admins have all permissions automatically.
func (u *User) Can(permission string) bool {
	if u.IsAdmin() {
		return true
	}
	for _, p := range u.Permissions {
		if p == permission || p == "*" {
			return true
		}
	}
	return false
}

// ApiKey represents a programmatic access key.
type ApiKey struct {
	ID        string
	UserID    string
	Name      string
	Key       string // hashed
	Scopes    []string
	LastUsedAt *time.Time
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// WebAuthnCredential represents a Passkey/WebAuthn credential.
type WebAuthnCredential struct {
	ID              []byte
	UserID          string
	PublicKey       []byte
	AttestationType string
	Transport       []string
	SignCount       uint32
	CreatedAt       time.Time
}

// Session represents a stateful user session.
type Session struct {
	ID        string
	UserID    string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// RefreshToken represents a persisted refresh token.
type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	IsRevoked bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

// IsExpired returns true if the refresh token has passed its expiry.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid returns true if the token is neither revoked nor expired.
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && !rt.IsExpired()
}

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
)

// Notification is the core notification entity.
type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Title     string
	Message   string
	Data      map[string]interface{}
	IsRead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
