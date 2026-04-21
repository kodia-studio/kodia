package providers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
)

type UserProvider struct{}

func NewUserProvider() *UserProvider {
	return &UserProvider{}
}

func (p *UserProvider) Name() string {
	return "kodia:user"
}

func (p *UserProvider) Register(app *kodia.App) error {
	// 1. Repositories
	userRepo := postgres.NewUserRepository(app.DB)

	// Auto-Migrate user models (Professional Framework standard)
	if err := postgres.AutoMigrate(app.DB); err != nil {
		app.Log.Error("Failed to auto-migrate user models", zap.Error(err))
		return err
	}

	// 2. Services
	userService := services.NewUserService(userRepo, app.Log)
	app.Set("user_service", userService)

	// 3. Handlers
	validate := validator.New()
	userHandler := handlers.NewUserHandler(userService, validate, app.Log)
	app.Set("user_handler", userHandler)

	return nil
}

func (p *UserProvider) Boot(app *kodia.App) error {
	if app.Router != nil {
		p.registerRoutes(app)
	}
	return nil
}

func (p *UserProvider) registerRoutes(app *kodia.App) {
	userHandler := app.MustGet("user_handler").(*handlers.UserHandler)
	jwtManager := app.MustGet("jwt_manager").(*jwt.Manager)
	
	api := app.Router.Group("/api")
	users := api.Group("/users")
	users.Use(middleware.Auth(jwtManager))
	{
		users.GET("/me", userHandler.GetMe)
		users.POST("/me/change-password", userHandler.ChangePassword)

		admin := users.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("", userHandler.GetAll)
			admin.GET("/:id", userHandler.GetByID)
			admin.PATCH("/:id", userHandler.Update)
			admin.DELETE("/:id", userHandler.Delete)
		}
	}
}
