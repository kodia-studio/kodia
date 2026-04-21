package authsocial

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/pkg/kodia"
	"gorm.io/gorm"
)

// SocialAccount represents a user's social login account link.
type SocialAccount struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	UserID     string    `gorm:"type:uuid;index"`
	Provider   string    `gorm:"type:varchar(50);index"`
	ProviderID string    `gorm:"type:varchar(255);index"`
	Email      string    `gorm:"type:varchar(255)"`
	Name       string    `gorm:"type:varchar(255)"`
	AvatarURL  string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli"`
}

// TableName specifies the database table name.
func (SocialAccount) TableName() string {
	return "social_accounts"
}

// SocialAccountRepository defines the interface for social account persistence.
type SocialAccountRepository interface {
	FindByProvider(ctx context.Context, provider, providerID string) (*SocialAccount, error)
	Create(ctx context.Context, account *SocialAccount) error
	FindByUserIDAndProvider(ctx context.Context, userID, provider string) (*SocialAccount, error)
}

// socialAccountRepository implements SocialAccountRepository using GORM.
type socialAccountRepository struct {
	db *gorm.DB
}

// NewSocialAccountRepository creates a new social account repository.
func NewSocialAccountRepository(db *gorm.DB) SocialAccountRepository {
	return &socialAccountRepository{db: db}
}

// FindByProvider finds a social account by provider and provider ID.
func (r *socialAccountRepository) FindByProvider(ctx context.Context, provider, providerID string) (*SocialAccount, error) {
	var account SocialAccount
	result := r.db.WithContext(ctx).
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&account)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error, return nil
		}
		return nil, result.Error
	}

	return &account, nil
}

// Create creates a new social account.
func (r *socialAccountRepository) Create(ctx context.Context, account *SocialAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// FindByUserIDAndProvider finds a social account by user ID and provider.
func (r *socialAccountRepository) FindByUserIDAndProvider(ctx context.Context, userID, provider string) (*SocialAccount, error) {
	var account SocialAccount
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&account)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &account, nil
}

// GetDB extracts database connection from app for use in other files.
func GetDB(app *kodia.App) *gorm.DB {
	return app.DB
}
