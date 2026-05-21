package pulse

import (
	"context"

	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ServiceProvider struct {
	manager *Manager
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (p *ServiceProvider) Name() string {
	return "kodia:pulse"
}

func (p *ServiceProvider) Register(app *kodia.App) error {
	p.manager = NewManager(app.Log)
	app.Set("pulse_manager", p.manager)

	app.Log = app.Log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		pulseCore := NewCore(p.manager, zapcore.WarnLevel)
		return zapcore.NewTee(core, pulseCore)
	}))

	return nil
}

func (p *ServiceProvider) Boot(app *kodia.App) error {
	ctx := context.Background()
	go p.manager.Run(ctx)

	if app.Router != nil {
		jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")
		handler := NewHandler(p.manager, app.Log)

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
