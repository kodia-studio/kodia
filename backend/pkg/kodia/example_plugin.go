package kodia

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/**
 * Example Kodia Plugin 🐨🚀
 * This plugin demonstrates how to use the new Hook system and formal Metadata.
 */

type AuditLoggerPlugin struct {
	// Any dependencies or configuration
}

func NewAuditLoggerPlugin() *AuditLoggerPlugin {
	return &AuditLoggerPlugin{}
}

// Name returns the unique name of the provider.
func (p *AuditLoggerPlugin) Name() string {
	return "AuditLoggerPlugin"
}

// Metadata provides formal information about the plugin.
func (p *AuditLoggerPlugin) Metadata() PluginMetadata {
	return PluginMetadata{
		ID:          "io.kodia.audit-logger",
		Name:        "Institutional Audit Logger",
		Version:     "1.0.0",
		Author:      "Kodia Core Team",
		Description: "Automatically logs institutional events via hooks.",
	}
}

// Register binds services to the container.
func (p *AuditLoggerPlugin) Register(app *App) error {
	app.Log.Info("Audit Logger Plugin registered")
	return nil
}

// Boot sets up listeners and routes.
func (p *AuditLoggerPlugin) Boot(app *App) error {
	// Listen to a core hook
	app.Hooks.Listen("user.login", func(data any) {
		email := data.(string)
		app.Log.Info("AUDIT: User logged in", zap.String("email", email))
	})

	// Dispatch a plugin-specific hook
	app.Hooks.Dispatch("audit.started", p.Metadata().ID)
	
	return nil
}

// RegisterRoutes adds plugin-specific endpoints.
func (p *AuditLoggerPlugin) RegisterRoutes(router *gin.Engine, app *App) error {
	router.GET("/api/audit/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "active",
			"plugin": p.Metadata().Name,
		})
	})
	return nil
}
