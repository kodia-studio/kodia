package seeders

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/hash"
	"gorm.io/gorm"
)

// UserSeeder seeds the database with fake user records.
type UserSeeder struct{}

// Run executes the seeding logic for users.
func (s *UserSeeder) Run(db *gorm.DB) error {
	// Create 20 fake users
	users := make([]interface{}, 20)

	for i := 0; i < 20; i++ {
		password := "TestPassword123!"
		hashedPassword, err := hash.Make(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user := &domain.User{
			ID:       uuid.New().String(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: hashedPassword,
			Role:     domain.RoleUser,
			IsActive: true,
		}

		users[i] = user
	}

	// Insert all users
	if err := db.CreateInBatches(users, 10).Error; err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}

	return nil
}
