package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migration defines the interface for versioned database migrations.
type Migration interface {
	Up(schema *Schema) error
	Down(schema *Schema) error
}

// MigrationRecord tracks which migrations have been executed.
type MigrationRecord struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`
	Name  string    `gorm:"uniqueIndex;size:255;not null"`
	Batch int       `gorm:"not null"`
	RunAt time.Time `gorm:"not null"`
}

// TableName specifies the table name for MigrationRecord.
func (MigrationRecord) TableName() string {
	return "kodia_migrations"
}

// Migrator manages versioned database migrations with atomic execution and tracking.
type Migrator struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(db *gorm.DB, log *zap.Logger) *Migrator {
	return &Migrator{
		db:  db,
		log: log,
	}
}

// EnsureTable creates the migrations tracking table if it doesn't exist.
func (m *Migrator) EnsureTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// Run executes a migration if it hasn't been run before.
// Returns nil if migration already executed or if it succeeds.
// Returns error if migration execution fails.
func (m *Migrator) Run(name string, migration Migration) error {
	// Check if migration already ran
	var count int64
	if err := m.db.Model(&MigrationRecord{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}
	if count > 0 {
		m.log.Debug("Migration already executed", zap.String("name", name))
		return nil
	}

	// Get the next batch number
	var lastBatch int
	if err := m.db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error; err != nil {
		return fmt.Errorf("failed to get last batch number: %w", err)
	}
	nextBatch := lastBatch + 1

	// Execute the migration
	m.log.Info("Running migration", zap.String("name", name), zap.Int("batch", nextBatch))
	schema := NewSchema(m.db)
	if err := migration.Up(schema); err != nil {
		return fmt.Errorf("migration %s failed: %w", name, err)
	}

	// Record the migration
	record := MigrationRecord{
		Name:  name,
		Batch: nextBatch,
		RunAt: time.Now(),
	}
	if err := m.db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to record migration %s: %w", name, err)
	}

	m.log.Info("Migration completed", zap.String("name", name))
	return nil
}

// RunAll executes all migrations in order, skipping those already executed.
func (m *Migrator) RunAll(migrations []struct {
	Name      string
	Migration Migration
}) error {
	for _, entry := range migrations {
		if err := m.Run(entry.Name, entry.Migration); err != nil {
			return err
		}
	}
	return nil
}
