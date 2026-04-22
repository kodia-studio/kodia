package providers

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/infrastructure/static"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

// StaticProvider manages serving embedded frontend assets.
type StaticProvider struct{}

func NewStaticProvider() *StaticProvider {
	return &StaticProvider{}
}

func (p *StaticProvider) Name() string {
	return "kodia:static"
}

func (p *StaticProvider) Register(app *kodia.App) error {
	// Only serve static files if we are in production or if configured
	if app.Config.IsProduction() {
		app.Log.Info("Frontend embedding enabled (Production Mode)")
		
		// Create subfilesystem for the dist folder
		sub, err := fs.Sub(static.DistFS, "dist")
		if err != nil {
			return err
		}

		staticServer := http.FS(sub)

		// Serve all static assets
		app.Router.NoRoute(func(c *gin.Context) {
			// If request is for API, don't serve static file (handles 404 naturally)
			if c.Request.URL.Path[:4] == "/api" {
				c.JSON(404, gin.H{"error": "route not found"})
				return
			}
			
			// Serve static file
			http.FileServer(staticServer).ServeHTTP(c.Writer, c.Request)
		})
	}
	return nil
}

func (p *StaticProvider) Boot(app *kodia.App) error {
	return nil
}
