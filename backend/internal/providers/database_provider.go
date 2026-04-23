package providers

import (
	"context"

	"github.com/kodia-studio/kodia/internal/infrastructure/database"
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
	db, err := database.New(app.Config, app.Log)
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
	return nil
}
