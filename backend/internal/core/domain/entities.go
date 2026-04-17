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
	IsActive  bool
	AvatarURL *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// IsAdmin returns true if the user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
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
