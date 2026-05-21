package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/validation"
	"github.com/kodia-studio/kodia/pkg/version"
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
	validate := validation.New()
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
	userHandler := kodia.MustResolve[*handlers.UserHandler](app, "user_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	api := app.Router.Group("/api")

	// /api/v1/users/... — versioned routes (canonical)
	v1 := api.Group("/v1")
	v1.Use(version.Middleware())
	registerUserRoutes(v1.Group("/users"), userHandler, jwtManager)

	// /api/users/... — unversioned legacy routes (backward compatibility)
	registerUserRoutes(api.Group("/users"), userHandler, jwtManager)
}

func registerUserRoutes(g *gin.RouterGroup, h *handlers.UserHandler, jwtManager *jwt.Manager) {
	g.Use(middleware.Auth(jwtManager))
	{
		g.GET("/me", h.GetMe)
		g.POST("/me/change-password", h.ChangePassword)
		g.DELETE("/me", h.DeleteMe)

		admin := g.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("", h.Create)
			admin.GET("", h.GetAll)
			admin.GET("/:id", h.GetByID)
			admin.PATCH("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
			admin.GET("/trashed", h.GetTrashed)
			admin.POST("/:id/restore", h.Restore)
			admin.POST("/:id/force-delete", h.ForceDelete)
		}
	}
}
