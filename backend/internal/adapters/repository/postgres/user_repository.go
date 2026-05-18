// Package postgres contains GORM-based repository implementations for Kodia Framework.
// These implementations satisfy the ports.UserRepository and ports.RefreshTokenRepository interfaces.
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/database"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"github.com/lib/pq"
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
	IsVerified bool       `gorm:"column:is_verified;not null;default:false"`

	// 2FA Security
	TwoFactorEnabled      bool           `gorm:"column:two_factor_enabled;not null;default:false"`
	TwoFactorSecret       string         `gorm:"column:two_factor_secret"`
	TwoFactorRecoveryCodes pq.StringArray `gorm:"column:two_factor_recovery_codes;type:text[]"`

	AvatarURL *string        `gorm:"column:avatar_url"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
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
		IsVerified: g.IsVerified,
		TwoFactorEnabled: g.TwoFactorEnabled,
		TwoFactorSecret: g.TwoFactorSecret,
		TwoFactorRecoveryCodes: []string(g.TwoFactorRecoveryCodes),
		AvatarURL: g.AvatarURL,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: database.FromDeletedAt(g.DeletedAt),
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
		IsVerified: u.IsVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		TwoFactorSecret: u.TwoFactorSecret,
		TwoFactorRecoveryCodes: pq.StringArray(u.TwoFactorRecoveryCodes),
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: database.ToDeletedAt(u.DeletedAt),
	}
}

// UserRepository is the GORM implementation of ports.UserRepository.
type UserRepository struct {
	database.BaseRepository[gormUser]
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: database.NewBaseRepository[gormUser](db)}
}

// AutoMigrate runs the GORM auto-migration for the user and auth models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&gormUser{}, &gormRefreshToken{},
		&gormProduct{}, &gormProductVariant{},
		&gormOrder{}, &gormOrderItem{}, &gormPayment{},
		&gormCoupon{},
	)
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	m := fromDomainUser(user)
	result := r.DB().WithContext(ctx).Create(m)
	return result.Error
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var m gormUser
	result := r.DB().WithContext(ctx).Where("id = ?", id).First(&m)
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
	result := r.DB().WithContext(ctx).Where("email = ?", email).First(&m)
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
	total, err := r.Query().FindPaginated(ctx, &models, params)
	if err != nil {
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
	return r.DB().WithContext(ctx).Save(m).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.RawDelete(ctx, id)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.Query().Where("email = ?", email).Exists(ctx)
}

func (r *UserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	return r.Query().Where("role = ?", role).Count(ctx)
}

func (r *UserRepository) Restore(ctx context.Context, id string) error {
	return r.RawRestore(ctx, id)
}

func (r *UserRepository) ForceDelete(ctx context.Context, id string) error {
	return r.RawForceDelete(ctx, id)
}

func (r *UserRepository) FindTrashed(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error) {
	var models []gormUser
	total, err := r.Query().OnlyTrashed().FindPaginated(ctx, &models, params)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		mc := m
		users[i] = mc.toDomain()
	}
	return users, total, nil
}
