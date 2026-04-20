package kodia

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServiceProvider defines the contract for a Kodia module.
type ServiceProvider interface {
	// Name returns the unique name of the provider.
	Name() string

	// Register is called early in the boot process to bind things to the container.
	Register(app *App) error

	// Boot is called after all providers are registered.
	// This is where routes and event listeners should be defined.
	Boot(app *App) error
}

// RouterProvider is an optional interface for providers that need to register routes.
type RouterProvider interface {
	RegisterRoutes(router *gin.Engine, app *App) error
}

// AppContext provides access to core application components.
type AppContext struct {
	Config *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Router *gin.Engine
}
