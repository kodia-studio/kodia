package providers

import (
	"context"
	"fmt"

	infradb "github.com/kodia-studio/kodia/internal/infrastructure/database"
	migrations "github.com/kodia-studio/kodia/internal/infrastructure/database/migrations/go"
	"github.com/kodia-studio/kodia/pkg/database"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

type DatabaseProvider struct{}

func NewDatabaseProvider() *DatabaseProvider {
	return &DatabaseProvider{}
}

func (p *DatabaseProvider) Name() string {
	return "kodia:database"
}

func (p *DatabaseProvider) Register(app *kodia.App) error {
	db, err := infradb.New(app.Config, app.Log)
	if err != nil {
		return err
	}
	app.DB = db

	// Register DB cleanup
	app.RegisterCleanupTask(func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		app.Log.Info("Closing database connection...")
		return sqlDB.Close()
	})

	return nil
}

func (p *DatabaseProvider) Boot(app *kodia.App) error {
	migrator := database.NewMigrator(app.DB, app.Log)

	// Ensure the migrations tracking table exists
	if err := migrator.EnsureTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	// Run all registered migrations
	for _, entry := range migrations.All() {
		if err := migrator.Run(entry.Name, entry.Migration.(database.Migration)); err != nil {
			return err
		}
	}

	return nil
}
