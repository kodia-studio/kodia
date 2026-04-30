package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/database"
	"gorm.io/gorm"
)

// gormRole is the GORM model for the roles table.
type gormRole struct {
	ID          string           `gorm:"column:id;primaryKey"`
	Name        string           `gorm:"column:name;uniqueIndex;not null"`
	Description string           `gorm:"column:description"`
	Permissions []gormPermission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt   `gorm:"column:deleted_at;index"`
}

func (gormRole) TableName() string { return "roles" }

// gormPermission is the GORM model for the permissions table.
type gormPermission struct {
	ID          string         `gorm:"column:id;primaryKey"`
	Name        string         `gorm:"column:name;uniqueIndex;not null"`
	Description string         `gorm:"column:description"`
	Group       string         `gorm:"column:group;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormPermission) TableName() string { return "permissions" }

// gormUserRole is the join table for user_roles.
type gormUserRole struct {
	UserID string `gorm:"column:user_id;primaryKey"`
	RoleID string `gorm:"column:role_id;primaryKey"`
}

func (gormUserRole) TableName() string { return "user_roles" }

// toDomain converts gormRole to domain.RoleEntity.
func (g *gormRole) toDomain() *domain.RoleEntity {
	permissions := make([]string, len(g.Permissions))
	for i, p := range g.Permissions {
		permissions[i] = p.Name
	}

	return &domain.RoleEntity{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Permissions: permissions,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
		DeletedAt:   database.FromDeletedAt(g.DeletedAt),
	}
}

// fromDomainRole converts domain.RoleEntity to gormRole.
func fromDomainRole(r *domain.RoleEntity) *gormRole {
	return &gormRole{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   database.ToDeletedAt(r.DeletedAt),
	}
}

// toDomain converts gormPermission to domain.PermissionEntity.
func (g *gormPermission) toDomain() *domain.PermissionEntity {
	return &domain.PermissionEntity{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Group:       g.Group,
		CreatedAt:   g.CreatedAt,
		DeletedAt:   database.FromDeletedAt(g.DeletedAt),
	}
}

// fromDomainPermission converts domain.PermissionEntity to gormPermission.
func fromDomainPermission(p *domain.PermissionEntity) *gormPermission {
	return &gormPermission{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Group:       p.Group,
		CreatedAt:   p.CreatedAt,
		DeletedAt:   database.ToDeletedAt(p.DeletedAt),
	}
}

// RoleRepository is the GORM implementation of ports.RoleRepository.
type RoleRepository struct {
	database.BaseRepository[gormRole]
}

// NewRoleRepository creates a new GORM-backed RoleRepository.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{BaseRepository: database.NewBaseRepository[gormRole](db)}
}

// PermissionRepository is the GORM implementation of ports.PermissionRepository.
type PermissionRepository struct {
	database.BaseRepository[gormPermission]
}

// NewPermissionRepository creates a new GORM-backed PermissionRepository.
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{BaseRepository: database.NewBaseRepository[gormPermission](db)}
}

// AutoMigrateRoles runs GORM auto-migration for roles and permissions tables.
func AutoMigrateRoles(db *gorm.DB) error {
	return db.AutoMigrate(
		&gormRole{},
		&gormPermission{},
		&gormUserRole{},
	)
}

// --- RoleRepository Implementation ---

func (r *RoleRepository) Create(ctx context.Context, role *domain.RoleEntity) error {
	m := fromDomainRole(role)

	// Load permissions if provided
	if len(role.Permissions) > 0 {
		var perms []gormPermission
		if err := r.DB().WithContext(ctx).Where("name IN ?", role.Permissions).Find(&perms).Error; err != nil {
			return err
		}
		m.Permissions = perms
	}

	result := r.DB().WithContext(ctx).Create(m)
	return result.Error
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*domain.RoleEntity, error) {
	var m gormRole
	err := r.Query().Preload("Permissions").Where("name = ?", name).First(ctx, &m)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m.toDomain(), nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]*domain.RoleEntity, error) {
	var roles []gormRole
	if err := r.Query().Preload("Permissions").Find(ctx, &roles); err != nil {
		return nil, err
	}

	entities := make([]*domain.RoleEntity, len(roles))
	for i, role := range roles {
		entities[i] = role.toDomain()
	}
	return entities, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.RoleEntity) error {
	m := fromDomainRole(role)

	// Update basic fields
	if err := r.DB().WithContext(ctx).Model(m).Updates(m).Error; err != nil {
		return err
	}

	// Update permissions association if provided
	if len(role.Permissions) > 0 {
		var perms []gormPermission
		if err := r.DB().WithContext(ctx).Where("name IN ?", role.Permissions).Find(&perms).Error; err != nil {
			return err
		}
		if err := r.DB().WithContext(ctx).Model(m).Association("Permissions").Replace(perms); err != nil {
			return err
		}
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.RawDelete(ctx, id)
}

func (r *RoleRepository) Restore(ctx context.Context, id string) error {
	return r.RawRestore(ctx, id)
}

func (r *RoleRepository) ForceDelete(ctx context.Context, id string) error {
	return r.RawForceDelete(ctx, id)
}

func (r *RoleRepository) AssignToUser(ctx context.Context, userID, roleName string) error {
	// Find the role ID
	var role gormRole
	if err := r.DB().WithContext(ctx).Select("id").Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	// Create association in user_roles
	return r.DB().WithContext(ctx).Create(&gormUserRole{
		UserID: userID,
		RoleID: role.ID,
	}).Error
}

func (r *RoleRepository) RevokeFromUser(ctx context.Context, userID, roleName string) error {
	// Find the role ID
	var role gormRole
	if err := r.DB().WithContext(ctx).Select("id").Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	// Delete the association
	return r.DB().WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, role.ID).Delete(&gormUserRole{}).Error
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roles []gormRole
	err := r.Query().
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(ctx, &roles)

	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}

// --- PermissionRepository Implementation ---

func (p *PermissionRepository) FindAll(ctx context.Context) ([]*domain.PermissionEntity, error) {
	var perms []gormPermission
	if err := p.Query().Find(ctx, &perms); err != nil {
		return nil, err
	}

	entities := make([]*domain.PermissionEntity, len(perms))
	for i, perm := range perms {
		entities[i] = perm.toDomain()
	}
	return entities, nil
}

func (p *PermissionRepository) FindByName(ctx context.Context, name string) (*domain.PermissionEntity, error) {
	var m gormPermission
	err := p.Query().Where("name = ?", name).First(ctx, &m)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m.toDomain(), nil
}

func (p *PermissionRepository) Create(ctx context.Context, perm *domain.PermissionEntity) error {
	m := fromDomainPermission(perm)
	return p.DB().WithContext(ctx).Create(m).Error
}

func (p *PermissionRepository) FindByGroup(ctx context.Context, group string) ([]*domain.PermissionEntity, error) {
	var perms []gormPermission
	if err := p.Query().Where("group = ?", group).Find(ctx, &perms); err != nil {
		return nil, err
	}

	entities := make([]*domain.PermissionEntity, len(perms))
	for i, perm := range perms {
		entities[i] = perm.toDomain()
	}
	return entities, nil
}
