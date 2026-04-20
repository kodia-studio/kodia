package search

import (
	"github.com/hibiken/asynq"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"gorm.io/gorm"
)

type SearchServiceProvider struct {
	manager *SearchManager
}

func NewServiceProvider() *SearchServiceProvider {
	return &SearchServiceProvider{}
}

func (p *SearchServiceProvider) Name() string {
	return "kodia:search"
}

func (p *SearchServiceProvider) Register(app *kodia.App) error {
	// Initialize Manager with Meilisearch as default driver
	// In production, these should come from app.Config
	host := "http://localhost:7700" 
	apiKey := "" // Default Meilisearch dev key

	manager := NewSearchManager("meilisearch")
	driver := NewMeiliSearchDriver(host, apiKey)
	manager.RegisterDriver("meilisearch", driver)

	p.manager = manager
	app.Set("search", manager)

	return nil
}

func (p *SearchServiceProvider) Boot(app *kodia.App) error {
	// Register Asynq handlers
	if _, ok := app.Get("asynq_server"); ok {
		// This registration pattern depends on how your asynq_server is set up
		// Assuming a standard mux registration here
		// app.Log.Info("Registering search background jobs")
	}

	// Register GORM hooks for automatic indexing
	if dbRaw, ok := app.Get("db"); ok {
		if db, ok := dbRaw.(*gorm.DB); ok {
			p.registerGormHooks(db, app)
		}
	}

	return nil
}

func (p *SearchServiceProvider) registerGormHooks(db *gorm.DB, app *kodia.App) {
	asynqClient := app.MustGet("asynq_client").(*asynq.Client)

	// After Create/Update Hook
	db.Callback().Create().After("gorm:create").Register("kodia:search:index", func(d *gorm.DB) {
		if searchable, ok := d.Statement.Model.(Searchable); ok {
			DispatchIndexTask(asynqClient, searchable.SearchIndex(), searchable.SearchID(), searchable.ToSearchMap())
		}
	})

	db.Callback().Update().After("gorm:update").Register("kodia:search:sync", func(d *gorm.DB) {
		if searchable, ok := d.Statement.Model.(Searchable); ok {
			DispatchIndexTask(asynqClient, searchable.SearchIndex(), searchable.SearchID(), searchable.ToSearchMap())
		}
	})

	// After Delete Hook
	db.Callback().Delete().After("gorm:delete").Register("kodia:search:delete", func(d *gorm.DB) {
		if searchable, ok := d.Statement.Model.(Searchable); ok {
			DispatchDeleteTask(asynqClient, searchable.SearchIndex(), searchable.SearchID())
		}
	})
}
