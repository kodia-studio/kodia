package kodia

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	// Hooks manages application-wide event listeners
	Hooks *HookManager

	// List of tasks to execute during graceful shutdown
	cleanupTasks []func(context.Context) error
}

func NewApp(cfg *config.Config, log *zap.Logger) *App {
	return &App{
		Config:    cfg,
		Log:       log,
		container: make(map[string]interface{}),
		Hooks:     NewHookManager(),
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

// RegisterCleanupTask adds a function to be called during graceful shutdown.
func (a *App) RegisterCleanupTask(task func(context.Context) error) {
	a.cleanupTasks = append(a.cleanupTasks, task)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Shutdown HTTP Server first to stop receiving new requests
	if err := server.Shutdown(ctx); err != nil {
		a.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	// 2. Execute all registered cleanup tasks (DB, Redis, Tracing, etc.)
	for i := len(a.cleanupTasks) - 1; i >= 0; i-- {
		if err := a.cleanupTasks[i](ctx); err != nil {
			a.Log.Error("Cleanup task failed", zap.Error(err))
		}
	}

	a.Log.Info("Application stopped gracefully")
	return nil
}

// NewTestServer starts the app in a test environment and returns a test server.
func (a *App) NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	
	// Boot the app (RegisterProviders must be called before NewTestServer)
	if err := a.Boot(); err != nil {
		t.Fatalf("failed to boot app for testing: %v", err)
	}
	
	// Use the router as the handler
	ts := httptest.NewServer(a.Router)
	
	// Register cleanup to shutdown the app and stop the server
	t.Cleanup(func() {
		ts.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		for i := len(a.cleanupTasks) - 1; i >= 0; i-- {
			_ = a.cleanupTasks[i](ctx)
		}
	})
	
	return ts
}

