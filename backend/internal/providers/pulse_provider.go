package providers

import (
	"context"

	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/observability"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PulseProvider struct {
	manager *observability.PulseManager
}

func NewPulseProvider() *PulseProvider {
	return &PulseProvider{}
}

func (p *PulseProvider) Name() string {
	return "kodia:pulse"
}

func (p *PulseProvider) Register(app *kodia.App) error {
	p.manager = observability.NewPulseManager(app.Log)
	app.Set("pulse_manager", p.manager)

	// Wrap the global logger with PulseCore to intercept Warn/Error logs
	app.Log = app.Log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		pulseCore := observability.NewPulseCore(p.manager, zapcore.WarnLevel)
		return zapcore.NewTee(core, pulseCore)
	}))

	return nil
}

func (p *PulseProvider) Boot(app *kodia.App) error {
	// 1. Start the Pulse Manager in a background goroutine
	ctx := context.Background()
	go p.manager.Run(ctx)

	// 2. Register real-time stream route if router is available
	if app.Router != nil {
		jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")
		handler := handlers.NewPulseHandler(p.manager, app.Log)

		pulse := app.Router.Group("/api/pulse")
		pulse.Use(middleware.Auth(jwtManager))
		pulse.Use(middleware.RequireRole("admin"))
		{
			pulse.GET("/stream", handler.Stream)
		}
	}

	app.Log.Info("Kodia Pulse telemetry engine booted and linked to router")
	return nil
}
