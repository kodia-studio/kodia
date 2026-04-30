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

// MigrationStatus represents the status of a migration.
type MigrationStatus struct {
	Name    string
	Batch   *int
	RunAt   *time.Time
	Pending bool
}

// Status returns the status of all registered migrations.
func (m *Migrator) Status(migrations []struct {
	Name      string
	Migration Migration
}) ([]MigrationStatus, error) {
	// Fetch all executed migrations
	var records []MigrationRecord
	if err := m.db.Order("batch DESC, id ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch migration records: %w", err)
	}

	// Build a map for quick lookup
	executedMap := make(map[string]*MigrationRecord)
	for i := range records {
		executedMap[records[i].Name] = &records[i]
	}

	// Build status list
	statuses := make([]MigrationStatus, 0, len(migrations))
	for _, entry := range migrations {
		status := MigrationStatus{
			Name:    entry.Name,
			Pending: true,
		}
		if record, exists := executedMap[entry.Name]; exists {
			status.Batch = &record.Batch
			status.RunAt = &record.RunAt
			status.Pending = false
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// Rollback rolls back the last N batches of migrations.
func (m *Migrator) Rollback(migrations []struct {
	Name      string
	Migration Migration
}, steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	// Build a map of migrations for quick lookup
	migrationMap := make(map[string]Migration)
	for _, entry := range migrations {
		migrationMap[entry.Name] = entry.Migration
	}

	// Get the maximum batch number
	var maxBatch int
	if err := m.db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch).Error; err != nil {
		return fmt.Errorf("failed to get max batch number: %w", err)
	}

	if maxBatch == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Rollback the last N batches
	for step := 0; step < steps; step++ {
		currentBatch := maxBatch - step

		// Get all migrations in the current batch (in reverse order)
		var records []MigrationRecord
		if err := m.db.Where("batch = ?", currentBatch).Order("id DESC").Find(&records).Error; err != nil {
			return fmt.Errorf("failed to fetch batch %d migrations: %w", currentBatch, err)
		}

		if len(records) == 0 {
			break
		}

		// Roll back each migration in reverse order
		for _, record := range records {
			migration, ok := migrationMap[record.Name]
			if !ok {
				return fmt.Errorf("migration %s not found in registry", record.Name)
			}

			m.log.Info("Rolling back migration", zap.String("name", record.Name), zap.Int("batch", record.Batch))
			schema := NewSchema(m.db)
			if err := migration.Down(schema); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", record.Name, err)
			}

			// Delete the migration record
			if err := m.db.Delete(&record).Error; err != nil {
				return fmt.Errorf("failed to delete migration record %s: %w", record.Name, err)
			}

			m.log.Info("Migration rolled back", zap.String("name", record.Name))
		}
	}

	return nil
}

// RollbackAll rolls back all migrations in reverse order of execution.
func (m *Migrator) RollbackAll(migrations []struct {
	Name      string
	Migration Migration
}) error {
	// Get the maximum batch number
	var maxBatch int
	if err := m.db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch).Error; err != nil {
		return fmt.Errorf("failed to get max batch number: %w", err)
	}

	if maxBatch == 0 {
		return nil // No migrations to rollback
	}

	return m.Rollback(migrations, maxBatch)
}

// Fresh drops all tables and re-runs all migrations.
func (m *Migrator) Fresh(migrations []struct {
	Name      string
	Migration Migration
}) error {
	m.log.Info("Starting fresh migration (dropping all tables)")

	// Get all tables
	tables, err := m.db.Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Drop all tables except the migrations table
	for _, table := range tables {
		if table != "kodia_migrations" {
			m.log.Debug("Dropping table", zap.String("table", table))
			if err := m.db.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}
	}

	// Clear migration records
	if err := m.db.Exec("DELETE FROM kodia_migrations").Error; err != nil {
		return fmt.Errorf("failed to clear migration records: %w", err)
	}

	m.log.Info("All tables dropped, re-running migrations")

	// Re-run all migrations
	return m.RunAll(migrations)
}

// Reset rolls back all migrations and then re-runs them.
func (m *Migrator) Reset(migrations []struct {
	Name      string
	Migration Migration
}) error {
	m.log.Info("Starting reset (rolling back all and re-running)")

	// Rollback all
	if err := m.RollbackAll(migrations); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	// Re-run all
	if err := m.RunAll(migrations); err != nil {
		return fmt.Errorf("failed to re-run migrations: %w", err)
	}

	m.log.Info("Reset completed successfully")
	return nil
}
