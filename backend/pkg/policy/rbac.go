package policy

import (
	"fmt"
	"sync"
)

// Permission represents a permission name.
type Permission string

// Role represents a role with a set of permissions.
type Role struct {
	Name        string
	Permissions map[Permission]bool
}

// RBACEngine provides in-memory role-based access control.
// It's built from database at startup and used for fast permission checks.
type RBACEngine struct {
	roles map[string]*Role
	mu    sync.RWMutex
}

// NewRBACEngine creates a new RBAC engine.
func NewRBACEngine() *RBACEngine {
	return &RBACEngine{
		roles: make(map[string]*Role),
	}
}

// DefineRole creates or updates a role with the given permissions.
// This method is chainable.
func (e *RBACEngine) DefineRole(name string, permissions ...Permission) *RBACEngine {
	e.mu.Lock()
	defer e.mu.Unlock()

	perms := make(map[Permission]bool)
	for _, p := range permissions {
		perms[p] = true
	}

	e.roles[name] = &Role{
		Name:        name,
		Permissions: perms,
	}

	return e
}

// Grant adds a permission to a role.
func (e *RBACEngine) Grant(roleName string, permission Permission) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	role, ok := e.roles[roleName]
	if !ok {
		return fmt.Errorf("role %s not found", roleName)
	}

	role.Permissions[permission] = true
	return nil
}

// Revoke removes a permission from a role.
func (e *RBACEngine) Revoke(roleName string, permission Permission) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	role, ok := e.roles[roleName]
	if !ok {
		return fmt.Errorf("role %s not found", roleName)
	}

	delete(role.Permissions, permission)
	return nil
}

// HasPermission checks if a role has a specific permission.
func (e *RBACEngine) HasPermission(roleName string, permission Permission) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	role, ok := e.roles[roleName]
	if !ok {
		return false
	}

	if role.Permissions["*"] {
		return true
	}

	return role.Permissions[permission]
}

// Can checks if a user with the given roles has a specific permission.
// Returns true if ANY role has the permission.
func (e *RBACEngine) Can(userRoles []string, permission Permission) bool {
	for _, roleName := range userRoles {
		if e.HasPermission(roleName, permission) {
			return true
		}
	}
	return false
}

// GetRole returns a role by name.
func (e *RBACEngine) GetRole(name string) (*Role, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	role, ok := e.roles[name]
	return role, ok
}

// AllRoles returns all defined roles.
func (e *RBACEngine) AllRoles() []Role {
	e.mu.RLock()
	defer e.mu.RUnlock()

	roles := make([]Role, 0, len(e.roles))
	for _, role := range e.roles {
		r := *role
		roles = append(roles, r)
	}

	return roles
}

// Clear removes all roles from the engine.
func (e *RBACEngine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.roles = make(map[string]*Role)
}
