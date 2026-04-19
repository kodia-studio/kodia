package seeders

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"gorm.io/gorm"
)

// ProductSeeder seeds the database with fake product records.
type ProductSeeder struct{}

// Run executes the seeding logic for products.
func (s *ProductSeeder) Run(db *gorm.DB) error {
	// Create 50 fake products
	products := make([]interface{}, 50)

	for i := 0; i < 50; i++ {
		product := &domain.Product{
			ID:          uuid.New().String(),
			Name:        gofakeit.ProductName(),
			Description: gofakeit.Sentence(10),
			Price:       float64(gofakeit.Price(10, 500)),
		}

		products[i] = product
	}

	// Insert all products
	if err := db.CreateInBatches(products, 10).Error; err != nil {
		return fmt.Errorf("failed to create products: %w", err)
	}

	return nil
}
