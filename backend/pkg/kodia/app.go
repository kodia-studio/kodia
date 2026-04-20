package kodia

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App is the main application kernel that manages providers and lifecycle.
type App struct {
	Config    *config.Config
	Log       *zap.Logger
	DB        *gorm.DB
	Router    *gin.Engine
	providers []ServiceProvider
	
	// Map to store services/dependencies that can be shared across providers
	container map[string]interface{}
}

func NewApp(cfg *config.Config, log *zap.Logger) *App {
	return &App{
		Config:    cfg,
		Log:       log,
		container: make(map[string]interface{}),
	}
}

// RegisterProviders adds multiple service providers to the app.
func (a *App) RegisterProviders(providers ...ServiceProvider) error {
	for _, p := range providers {
		a.Log.Debug("Registering provider", zap.String("name", p.Name()))
		if err := p.Register(a); err != nil {
			return fmt.Errorf("failed to register provider %s: %w", p.Name(), err)
		}
		a.providers = append(a.providers, p)
	}
	return nil
}

// Boot initializes all registered providers.
func (a *App) Boot() error {
	for _, p := range a.providers {
		a.Log.Debug("Booting provider", zap.String("name", p.Name()))
		if err := p.Boot(a); err != nil {
			return fmt.Errorf("failed to boot provider %s: %w", p.Name(), err)
		}
	}
	return nil
}

// Set stores a dependency in the app container.
func (a *App) Set(key string, value interface{}) {
	a.container[key] = value
}

// Get retrieves a dependency from the app container.
func (a *App) Get(key string) (interface{}, bool) {
	val, ok := a.container[key]
	return val, ok
}

// MustGet retrieves a dependency or panics if it doesn't exist.
func (a *App) MustGet(key string) interface{} {
	val, ok := a.Get(key)
	if !ok {
		panic(fmt.Sprintf("dependency %s not found in container", key))
	}
	return val
}

// Run starts the HTTP server with graceful shutdown.
func (a *App) Run() error {
	if a.Router == nil {
		return fmt.Errorf("router not initialized")
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.App.Port),
		Handler: a.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	a.Log.Info("Application is running", zap.String("url", fmt.Sprintf("http://localhost:%d", a.Config.App.Port)))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Log.Info("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.Log.Info("Application stopped gracefully")
	return nil
}
