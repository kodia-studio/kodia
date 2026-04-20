package providers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kodia_http "github.com/kodia-studio/kodia/internal/adapters/http"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
)

type HttpProvider struct{}

func NewHttpProvider() *HttpProvider {
	return &HttpProvider{}
}

func (p *HttpProvider) Name() string {
	return "kodia:http"
}

func (p *HttpProvider) Register(app *kodia.App) error {
	if app.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	
	engine := gin.New()

	// Global middleware
	engine.Use(middleware.Recovery(app.Log))
	engine.Use(middleware.Logger(app.Log))

	// Tracing & Metrics Middleware
	if app.Config.Observability.TracingEnabled {
		engine.Use(middleware.Tracing(app.Config.Observability.ServiceName))
	}
	if app.Config.Observability.MetricsEnabled {
		engine.Use(middleware.Metrics())
	}

	// Validate CORS configuration
	if err := kodia_http.ValidateCORSConfig(app.Config, app.Log); err != nil {
		app.Log.Fatal("CORS configuration validation failed", zap.Error(err))
	}

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     app.Config.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Trace-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Trace-ID"},
		AllowCredentials: true,
	}
	engine.Use(cors.New(corsConfig))

	app.Router = engine
	return nil
}

func (p *HttpProvider) Boot(app *kodia.App) error {
	// Add global health check
	api := app.Router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": app.Config.App.Name})
	})
	
	return nil
}
