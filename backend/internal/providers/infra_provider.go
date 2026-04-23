package providers

import (
	"context"
	"strings"

	"github.com/kodia-studio/kodia/internal/infrastructure/cache"
	"github.com/kodia-studio/kodia/internal/infrastructure/events"
	"github.com/kodia-studio/kodia/internal/infrastructure/mailer"
	"github.com/kodia-studio/kodia/internal/infrastructure/storage"
	"github.com/kodia-studio/kodia/internal/infrastructure/worker"
	audit_infra "github.com/kodia-studio/kodia/internal/infrastructure/audit"
	core_events "github.com/kodia-studio/kodia/internal/core/events"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
)

type InfraProvider struct{}

func NewInfraProvider() *InfraProvider {
	return &InfraProvider{}
}

func (p *InfraProvider) Name() string {
	return "kodia:infra"
}

func (p *InfraProvider) Register(app *kodia.App) error {
	// 1. Storage
	var storageProvider ports.StorageProvider
	switch strings.ToLower(app.Config.Storage.Provider) {
	case "s3":
		sp, err := storage.NewS3StorageProvider(app.Config)
		if err != nil {
			return err
		}
		storageProvider = sp
	default:
		storageProvider = storage.NewLocalStorageProvider(app.Config)
	}
	app.Set("storage", storageProvider)

	// 2. Cache & Redis Cleanup
	cacheProvider, err := cache.New(app.Config, app.Log)
	if err != nil {
		app.Log.Warn("Cache initialization failed", zap.Error(err))
	} else {
		app.Set("cache", cacheProvider)
		
		// Register Redis cleanup if it's a RedisProvider
		if rp, ok := cacheProvider.(*cache.RedisProvider); ok {
			app.RegisterCleanupTask(func(ctx context.Context) error {
				app.Log.Info("Closing Redis connection...")
				return rp.GetClient().Close()
			})
		}
	}

	// 3. Mailer
	mailProvider, err := mailer.NewSMTPMailer(app.Config, app.Log)
	if err != nil {
		return err
	}
	app.Set("mailer", mailProvider)

	// 4. Events & Workers
	queueProvider := worker.NewAsynqProvider(app.Config, app.Log)
	dispatcher := events.NewDispatcher(queueProvider, app.Log)
	app.Set("events", dispatcher)

	// 5. Audit Logging
	auditManager := audit_infra.NewManager(app.Log)
	if app.DB != nil {
		auditManager.AddLogger(audit_infra.NewGormLogger(app.DB))
	}
	app.Set("audit", auditManager)

	return nil
}

func (p *InfraProvider) Boot(app *kodia.App) error {
	// Register event listeners
	if dispatcher, ok := app.Get("events"); ok {
		core_events.RegisterEvents(dispatcher.(ports.EventDispatcher))
	}
	return nil
}
