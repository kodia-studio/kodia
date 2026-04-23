package tests

import (
	"fmt"
	"testing"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Factory provides methods to create domain entities for testing
type Factory struct {
	db *gorm.DB
	t  *testing.T
}

// NewFactory creates a new test data factory
func NewFactory(t *testing.T, db *gorm.DB) *Factory {
	return &Factory{db: db, t: t}
}

// CreateUser creates a user with the given overrides
func (f *Factory) CreateUser(overrides ...func(*domain.User)) *domain.User {
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	user := &domain.User{
		Email:    fmt.Sprintf("user-%p@example.com", &overrides),
		Password: string(password),
		Role:     "user",
	}

	for _, override := range overrides {
		override(user)
	}

	if err := f.db.Create(user).Error; err != nil {
		f.t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

// CreateAdmin creates an admin user
func (f *Factory) CreateAdmin(overrides ...func(*domain.User)) *domain.User {
	overrides = append(overrides, func(u *domain.User) {
		u.Role = "admin"
	})
	return f.CreateUser(overrides...)
}

// CreateRefreshToken creates a refresh token for a user
func (f *Factory) CreateRefreshToken(userID string, overrides ...func(*domain.RefreshToken)) *domain.RefreshToken {
	token := &domain.RefreshToken{
		UserID:    userID,
		Token:     fmt.Sprintf("token-%p", &overrides),
		IsRevoked: false,
	}

	for _, override := range overrides {
		override(token)
	}

	if err := f.db.Create(token).Error; err != nil {
		f.t.Fatalf("failed to create test refresh token: %v", err)
	}

	return token
}
