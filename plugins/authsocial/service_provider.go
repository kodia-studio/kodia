package authsocial

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/authsocial"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"go.uber.org/zap"
)

// ServiceProvider registers the social login service with the framework.
type ServiceProvider struct{}

// NewServiceProvider creates a new social auth service provider.
func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

// Name returns the provider name.
func (p *ServiceProvider) Name() string {
	return "kodia:authsocial"
}

// Register registers social auth dependencies.
func (p *ServiceProvider) Register(app *kodia.App) error {
	// Read OAuth credentials from environment
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		googleRedirectURL = app.Config.App.BaseURL + "/api/auth/social/google/callback"
	}

	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURL := os.Getenv("GITHUB_REDIRECT_URL")
	if githubRedirectURL == "" {
		githubRedirectURL = app.Config.App.BaseURL + "/api/auth/social/github/callback"
	}

	// Validate that at least one provider is configured
	if (googleClientID == "" && githubClientID == "") {
		return fmt.Errorf("at least one of GOOGLE_CLIENT_ID or GITHUB_CLIENT_ID must be set")
	}

	// Create OAuth providers
	var googleProvider authsocial.Provider
	var githubProvider authsocial.Provider

	if googleClientID != "" {
		googleProvider = authsocial.NewGoogleProvider(authsocial.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		})
	}

	if githubClientID != "" {
		githubProvider = authsocial.NewGitHubProvider(authsocial.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  githubRedirectURL,
			Scopes:       []string{"user:email"},
		})
	}

	// Get dependencies from container
	jwtManager := app.MustGet("jwt_manager").(*jwt.Manager)
	cache := app.MustGet("cache").(ports.CacheProvider)

	// Create repositories
	socialRepo := NewSocialAccountRepository(app.DB)
	userRepo := postgres.NewUserRepository(app.DB)

	// Create service
	service := NewSocialAuthService(socialRepo, userRepo, cache, jwtManager, app.Config.App.BaseURL, app.Log)

	// Create handler
	handler := NewSocialHandler(service, app.Config.App.FrontendURL, googleProvider, githubProvider, app.Log)

	// Register in container
	app.Set("social_auth_service", service)
	app.Set("social_handler", handler)

	app.Log.Info("Social auth service registered",
		zap.Bool("google_enabled", googleClientID != ""),
		zap.Bool("github_enabled", githubClientID != ""),
	)

	return nil
}

// Boot registers routes.
func (p *ServiceProvider) Boot(app *kodia.App) error {
	if app.Router == nil {
		return nil
	}

	handler := app.MustGet("social_handler").(*SocialHandler)

	// Register routes
	social := app.Router.Group("/api/auth/social")
	{
		social.GET("/:provider/redirect", handler.Redirect)
		social.GET("/:provider/callback", handler.Callback)
	}

	app.Log.Info("Social auth routes registered")
	return nil
}
