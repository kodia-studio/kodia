package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// Registry contains all registered seeders.
var Registry = []Seeder{
	// --- Seeder Registration Start ---
	&UserSeeder{},
	&ProductSeeder{},
	// --- Seeder Registration End ---
}

// RunAll executes all registered seeders.
func RunAll(db *gorm.DB) error {
	for _, seeder := range Registry {
		fmt.Printf("Running seeder: %T\n", seeder)
		if err := seeder.Run(db); err != nil {
			return fmt.Errorf("failed running seeder %T: %w", seeder, err)
		}
	}
	return nil
}
