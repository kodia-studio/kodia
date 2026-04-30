package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/validation"
	"github.com/kodia-studio/kodia/pkg/version"
)

type AuthProvider struct{}

func NewAuthProvider() *AuthProvider {
	return &AuthProvider{}
}

func (p *AuthProvider) Name() string {
	return "kodia:auth"
}

func (p *AuthProvider) Register(app *kodia.App) error {
	// 1. Initialize JWT Manager
	jwtManager := jwt.NewManager(
		app.Config.JWT.AccessSecret,
		app.Config.JWT.RefreshSecret,
		app.Config.JWT.AccessExpiryHours,
		app.Config.JWT.RefreshExpiryDays,
	)
	app.Set("jwt_manager", jwtManager)

	// 2. Repositories
	userRepo := postgres.NewUserRepository(app.DB)
	refreshRepo := postgres.NewRefreshTokenRepository(app.DB)

	// Retrieve infra from container
	cacheProvider := kodia.MustResolve[ports.CacheProvider](app, "cache")
	mailProvider := kodia.MustResolve[ports.Mailer](app, "mailer")

	// 3. Services
	authService := services.NewAuthService(
		userRepo, 
		refreshRepo, 
		jwtManager, 
		cacheProvider, 
		mailProvider, 
		app.Config.App.BaseURL,
		app.Config.App.FrontendURL,
		app.Log,
	)
	app.Set("auth_service", authService)

	// 4. Handlers
	validate := validation.New()
	authHandler := handlers.NewAuthHandler(authService, validate, app.Log)
	app.Set("auth_handler", authHandler)

	return nil
}

func (p *AuthProvider) Boot(app *kodia.App) error {
	// Register Routes
	if app.Router != nil {
		p.registerRoutes(app)
	}
	return nil
}

func (p *AuthProvider) registerRoutes(app *kodia.App) {
	authHandler := kodia.MustResolve[*handlers.AuthHandler](app, "auth_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	api := app.Router.Group("/api")

	// /api/v1/auth/... — versioned routes (canonical)
	v1 := api.Group("/v1")
	v1.Use(version.Middleware())
	registerAuthRoutes(v1.Group("/auth"), authHandler, jwtManager)

	// /api/auth/... — unversioned legacy routes (backward compatibility)
	registerAuthRoutes(api.Group("/auth"), authHandler, jwtManager)
}

func registerAuthRoutes(g *gin.RouterGroup, h *handlers.AuthHandler, jwtManager *jwt.Manager) {
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.RefreshToken)
	g.POST("/forgot-password", h.ForgotPassword)
	g.POST("/reset-password", h.ResetPassword)
	g.GET("/verify-email", h.VerifyEmail)
	g.POST("/2fa/login-verify", h.LoginVerify2FA)
	g.POST("/logout", middleware.Auth(jwtManager), h.Logout)

	protected := g.Group("")
	protected.Use(middleware.Auth(jwtManager))
	{
		protected.POST("/logout-all", h.LogoutAll)
		protected.GET("/me", h.Me)

		// 2FA Management
		mfa := protected.Group("/2fa")
		{
			mfa.POST("/enable", h.Enable2FA)
			mfa.POST("/verify", h.Verify2FA)
			mfa.DELETE("/disable", h.Disable2FA)
		}
	}
}
