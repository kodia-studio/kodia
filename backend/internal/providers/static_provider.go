package providers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

// distFS is the embedded file system containing the frontend build.
// Note: This matches the path used in 'kodia build'
//go:embed all:infrastructure/static/dist/*
var distFS embed.FS

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
		sub, err := fs.Sub(distFS, "infrastructure/static/dist")
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
