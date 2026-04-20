package providers

import (
	"strings"

	"github.com/kodia-studio/kodia/internal/infrastructure/cache"
	"github.com/kodia-studio/kodia/internal/infrastructure/events"
	"github.com/kodia-studio/kodia/internal/infrastructure/mailer"
	"github.com/kodia-studio/kodia/internal/infrastructure/storage"
	"github.com/kodia-studio/kodia/internal/infrastructure/worker"
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

	// 2. Cache
	cacheProvider, err := cache.New(app.Config, app.Log)
	if err != nil {
		app.Log.Warn("Cache initialization failed", zap.Error(err))
	} else {
		app.Set("cache", cacheProvider)
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

	return nil
}

func (p *InfraProvider) Boot(app *kodia.App) error {
	// Register event listeners
	if dispatcher, ok := app.Get("events"); ok {
		core_events.RegisterEvents(dispatcher.(ports.EventDispatcher))
	}
	return nil
}
