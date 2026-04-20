package providers

import (
	"context"

	"github.com/kodia-studio/kodia/internal/infrastructure/database"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/observability"
	"go.uber.org/zap"
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
	
	// Observability init
	obsManager := observability.NewManager(app.Config, app.Log)
	if err := obsManager.Init(context.Background()); err != nil {
		app.Log.Warn("Observability init failed", zap.Error(err))
	}
	app.Set("observability", obsManager)
	
	return nil
}

func (p *DatabaseProvider) Boot(app *kodia.App) error {
	return nil
}
