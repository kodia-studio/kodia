// Package postgres contains GORM-based repository implementations for Kodia Framework.
// These implementations satisfy the ports.UserRepository and ports.RefreshTokenRepository interfaces.
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"gorm.io/gorm"
)

// gormUser is the GORM model for the users table.
// It mirrors the domain.User entity but with GORM-specific tags.
// We keep this separate to avoid polluting the domain with framework concerns.
type gormUser struct {
	ID        string     `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name;not null"`
	Email     string     `gorm:"column:email;uniqueIndex;not null"`
	Password  string     `gorm:"column:password;not null"`
	Role      string     `gorm:"column:role;not null;default:'user'"`
	IsActive  bool       `gorm:"column:is_active;not null;default:true"`
	AvatarURL *string    `gorm:"column:avatar_url"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (gormUser) TableName() string { return "users" }

// toDomain converts a gormUser to a domain.User entity.
func (g *gormUser) toDomain() *domain.User {
	return &domain.User{
		ID:        g.ID,
		Name:      g.Name,
		Email:     g.Email,
		Password:  g.Password,
		Role:      domain.UserRole(g.Role),
		IsActive:  g.IsActive,
		AvatarURL: g.AvatarURL,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
	}
}

// fromDomain converts a domain.User to a gormUser.
func fromDomainUser(u *domain.User) *gormUser {
	return &gormUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// UserRepository is the GORM implementation of ports.UserRepository.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	m := fromDomainUser(user)
	result := r.db.WithContext(ctx).Create(m)
	return result.Error
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var m gormUser
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return m.toDomain(), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m gormUser
	result := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return m.toDomain(), nil
}

func (r *UserRepository) FindAll(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error) {
	var models []gormUser
	var total int64

	baseQuery := r.db.WithContext(ctx).Model(&gormUser{}).Where("deleted_at IS NULL")

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.Offset(params.Offset()).Limit(params.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		mc := m
		users[i] = mc.toDomain()
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	m := fromDomainUser(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&gormUser{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&gormUser{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	return count > 0, err
}
