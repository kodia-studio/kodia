package providers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
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

	// 3. Services
	authService := services.NewAuthService(userRepo, refreshRepo, jwtManager, app.Log)
	app.Set("auth_service", authService)

	// 4. Handlers
	validate := validator.New()
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
	authHandler := app.MustGet("auth_handler").(*handlers.AuthHandler)
	jwtManager := app.MustGet("jwt_manager").(*jwt.Manager)
	
	api := app.Router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.Auth(jwtManager), authHandler.Logout)

		protected := auth.Group("")
		protected.Use(middleware.Auth(jwtManager))
		{
			protected.POST("/logout-all", authHandler.LogoutAll)
			protected.GET("/me", authHandler.Me)
		}
	}
}
