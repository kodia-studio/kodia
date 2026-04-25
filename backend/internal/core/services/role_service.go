package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/policy"
	"go.uber.org/zap"
)

// RoleService implements ports.RoleService.
type RoleService struct {
	roleRepo       ports.RoleRepository
	permRepo       ports.PermissionRepository
	rbacEngine     *policy.RBACEngine
	log            *zap.Logger
}

// NewRoleService creates a new RoleService with its dependencies injected.
func NewRoleService(
	roleRepo ports.RoleRepository,
	permRepo ports.PermissionRepository,
	rbacEngine *policy.RBACEngine,
	log *zap.Logger,
) *RoleService {
	return &RoleService{
		roleRepo:   roleRepo,
		permRepo:   permRepo,
		rbacEngine: rbacEngine,
		log:        log,
	}
}

// CreateRole creates a new role with the given permissions.
func (s *RoleService) CreateRole(ctx context.Context, name, description string, permissions []string) (*domain.RoleEntity, error) {
	role := &domain.RoleEntity{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		s.log.Error("Failed to create role", zap.String("role_name", name), zap.Error(err))
		return nil, fmt.Errorf("create role: %w", err)
	}

	// Update RBAC engine
	perms := make([]policy.Permission, len(permissions))
	for i, p := range permissions {
		perms[i] = policy.Permission(p)
	}
	s.rbacEngine.DefineRole(name, perms...)

	s.log.Info("Role created", zap.String("role_name", name))
	return role, nil
}

// AssignRole assigns a role to a user.
func (s *RoleService) AssignRole(ctx context.Context, userID, roleName string) error {
	// Verify role exists
	_, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if err := s.roleRepo.AssignToUser(ctx, userID, roleName); err != nil {
		s.log.Error("Failed to assign role to user",
			zap.String("user_id", userID),
			zap.String("role_name", roleName),
			zap.Error(err),
		)
		return fmt.Errorf("assign role: %w", err)
	}

	s.log.Info("Role assigned to user",
		zap.String("user_id", userID),
		zap.String("role_name", roleName),
	)
	return nil
}

// RevokeRole revokes a role from a user.
func (s *RoleService) RevokeRole(ctx context.Context, userID, roleName string) error {
	if err := s.roleRepo.RevokeFromUser(ctx, userID, roleName); err != nil {
		s.log.Error("Failed to revoke role from user",
			zap.String("user_id", userID),
			zap.String("role_name", roleName),
			zap.Error(err),
		)
		return fmt.Errorf("revoke role: %w", err)
	}

	s.log.Info("Role revoked from user",
		zap.String("user_id", userID),
		zap.String("role_name", roleName),
	)
	return nil
}

// GetUserRoles retrieves all role names assigned to a user.
func (s *RoleService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		s.log.Error("Failed to get user roles",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return roles, nil
}

// GetAllRoles retrieves all defined roles.
func (s *RoleService) GetAllRoles(ctx context.Context) ([]*domain.RoleEntity, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		s.log.Error("Failed to get all roles", zap.Error(err))
		return nil, fmt.Errorf("get all roles: %w", err)
	}
	return roles, nil
}

// DeleteRole soft-deletes a role by its ID.
func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		s.log.Error("Failed to delete role",
			zap.String("role_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("delete role: %w", err)
	}

	s.log.Info("Role deleted", zap.String("role_id", id))
	return nil
}

// SyncEngineFromDB loads all roles and permissions from the database into the RBAC engine.
// This is called at startup and after role/permission changes.
func (s *RoleService) SyncEngineFromDB(ctx context.Context) error {
	// Clear existing engine state
	s.rbacEngine.Clear()

	// Load all roles with their permissions
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		s.log.Error("Failed to sync RBAC engine from database", zap.Error(err))
		return fmt.Errorf("sync rbac engine: %w", err)
	}

	// Define roles in the engine
	for _, role := range roles {
		perms := make([]policy.Permission, len(role.Permissions))
		for i, p := range role.Permissions {
			perms[i] = policy.Permission(p)
		}
		s.rbacEngine.DefineRole(role.Name, perms...)
	}

	s.log.Info("RBAC engine synced from database", zap.Int("roles", len(roles)))
	return nil
}
