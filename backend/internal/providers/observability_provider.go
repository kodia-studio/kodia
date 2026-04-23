package providers

import (
	"context"

	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/infrastructure/cache"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/observability"
	"github.com/redis/go-redis/v9"
)

type ObservabilityProvider struct{}

func NewObservabilityProvider() *ObservabilityProvider {
	return &ObservabilityProvider{}
}

func (p *ObservabilityProvider) Name() string {
	return "kodia:observability"
}

func (p *ObservabilityProvider) Register(app *kodia.App) error {
	// 1. Initialize Observability Manager (Tracing, Metrics, Sentry)
	obsManager := observability.NewManager(app.Config, app.Log)
	if err := obsManager.Init(context.Background()); err != nil {
		app.Log.Error("Observability Manager init failed")
	}
	app.Set("observability", obsManager)

	// Register cleanup for Observability
	app.RegisterCleanupTask(func(ctx context.Context) error {
		app.Log.Info("Shutting down observability stack...")
		obsManager.Shutdown(ctx)
		return nil
	})

	// 2. Initialize Health Handler
	var redisClient *redis.Client
	if cacheProv, ok := app.Get("cache"); ok {
		if rp, ok := cacheProv.(*cache.RedisProvider); ok {
			redisClient = rp.GetClient()
		}
	}

	healthHandler := handlers.NewHealthHandler(app.DB, redisClient, app.Log)
	app.Set("health_handler", healthHandler)

	return nil
}

func (p *ObservabilityProvider) Boot(app *kodia.App) error {
	if app.Router != nil {
		healthHandler := app.MustGet("health_handler").(*handlers.HealthHandler)
		
		// Register Health Routes
		api := app.Router.Group("/api/v1")
		{
			api.GET("/health", healthHandler.Ready)
			api.GET("/health/live", healthHandler.Live)
			api.GET("/health/ready", healthHandler.Ready)
		}
	}
	return nil
}
