package domain

import "time"

// RoleEntity represents a role with a set of permissions.
type RoleEntity struct {
	ID          string
	Name        string
	Description string
	Permissions []string  // permission names
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// PermissionEntity represents a specific permission in the system.
type PermissionEntity struct {
	ID          string
	Name        string
	Description string
	Group       string  // e.g. "users", "products", "admin"
	CreatedAt   time.Time
	DeletedAt   *time.Time
}
