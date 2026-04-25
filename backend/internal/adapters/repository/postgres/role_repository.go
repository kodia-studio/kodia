package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"gorm.io/gorm"
)

// gormRole is the GORM model for the roles table.
type gormRole struct {
	ID          string         `gorm:"column:id;primaryKey"`
	Name        string         `gorm:"column:name;uniqueIndex;not null"`
	Description string         `gorm:"column:description"`
	Permissions []gormPermission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time     `gorm:"column:deleted_at;index"`
}

func (gormRole) TableName() string { return "roles" }

// gormPermission is the GORM model for the permissions table.
type gormPermission struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;uniqueIndex;not null"`
	Description string    `gorm:"column:description"`
	Group     string    `gorm:"column:group;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
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
		DeletedAt:   g.DeletedAt,
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
		DeletedAt:   r.DeletedAt,
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
		DeletedAt:   g.DeletedAt,
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
		DeletedAt:   p.DeletedAt,
	}
}

// RoleRepository is the GORM implementation of ports.RoleRepository.
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new GORM-backed RoleRepository.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// PermissionRepository is the GORM implementation of ports.PermissionRepository.
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new GORM-backed PermissionRepository.
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
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
		if err := r.db.WithContext(ctx).Where("name IN ?", role.Permissions).Find(&perms).Error; err != nil {
			return err
		}
		m.Permissions = perms
	}

	result := r.db.WithContext(ctx).Create(m)
	return result.Error
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*domain.RoleEntity, error) {
	var m gormRole
	result := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.toDomain(), nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]*domain.RoleEntity, error) {
	var roles []gormRole
	result := r.db.WithContext(ctx).Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
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
	if err := r.db.WithContext(ctx).Model(m).Updates(m).Error; err != nil {
		return err
	}

	// Update permissions association if provided
	if len(role.Permissions) > 0 {
		var perms []gormPermission
		if err := r.db.WithContext(ctx).Where("name IN ?", role.Permissions).Find(&perms).Error; err != nil {
			return err
		}
		if err := r.db.WithContext(ctx).Model(m).Association("Permissions").Replace(perms); err != nil {
			return err
		}
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&gormRole{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *RoleRepository) AssignToUser(ctx context.Context, userID, roleName string) error {
	// Find the role ID
	var role gormRole
	if err := r.db.WithContext(ctx).Select("id").Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	// Create association in user_roles
	return r.db.WithContext(ctx).Create(&gormUserRole{
		UserID: userID,
		RoleID: role.ID,
	}).Error
}

func (r *RoleRepository) RevokeFromUser(ctx context.Context, userID, roleName string) error {
	// Find the role ID
	var role gormRole
	if err := r.db.WithContext(ctx).Select("id").Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	// Delete the association
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, role.ID).Delete(&gormUserRole{}).Error
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roles []gormRole
	result := r.db.WithContext(ctx).
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles)

	if result.Error != nil {
		return nil, result.Error
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
	result := p.db.WithContext(ctx).Find(&perms)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*domain.PermissionEntity, len(perms))
	for i, perm := range perms {
		entities[i] = perm.toDomain()
	}
	return entities, nil
}

func (p *PermissionRepository) FindByName(ctx context.Context, name string) (*domain.PermissionEntity, error) {
	var m gormPermission
	result := p.db.WithContext(ctx).Where("name = ?", name).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.toDomain(), nil
}

func (p *PermissionRepository) Create(ctx context.Context, perm *domain.PermissionEntity) error {
	m := fromDomainPermission(perm)
	return p.db.WithContext(ctx).Create(m).Error
}

func (p *PermissionRepository) FindByGroup(ctx context.Context, group string) ([]*domain.PermissionEntity, error) {
	var perms []gormPermission
	result := p.db.WithContext(ctx).Where("group = ?", group).Find(&perms)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*domain.PermissionEntity, len(perms))
	for i, perm := range perms {
		entities[i] = perm.toDomain()
	}
	return entities, nil
}
