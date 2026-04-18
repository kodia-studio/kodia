package seeders

import (
	"gorm.io/gorm"
)

// Seeder defines the interface for database seeding.
type Seeder interface {
	// Run executes the seeding logic.
	Run(db *gorm.DB) error
}
