package providers

import (
	"context"

	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/policy"
	"github.com/kodia-studio/kodia/pkg/validation"
	"go.uber.org/zap"
)

type RoleProvider struct{}

func NewRoleProvider() *RoleProvider {
	return &RoleProvider{}
}

func (p *RoleProvider) Name() string {
	return "kodia:role"
}

func (p *RoleProvider) Register(app *kodia.App) error {
	// 1. Repositories
	roleRepo := postgres.NewRoleRepository(app.DB)
	permRepo := postgres.NewPermissionRepository(app.DB)

	// Auto-Migrate role models
	if err := postgres.AutoMigrateRoles(app.DB); err != nil {
		app.Log.Error("Failed to auto-migrate role models", zap.Error(err))
		return err
	}

	// 2. RBAC Engine
	rbacEngine := policy.NewRBACEngine()
	app.Set("rbac", rbacEngine)

	// 3. Services
	roleService := services.NewRoleService(roleRepo, permRepo, rbacEngine, app.Log)
	app.Set("role_service", roleService)

	// Sync RBAC engine from database at startup
	if err := roleService.SyncEngineFromDB(context.Background()); err != nil {
		app.Log.Error("Failed to sync RBAC engine from database", zap.Error(err))
		return err
	}

	// 4. Handlers
	validate := validation.New()
	roleHandler := handlers.NewRoleHandler(roleService, validate, app.Log)
	app.Set("role_handler", roleHandler)

	return nil
}

func (p *RoleProvider) Boot(app *kodia.App) error {
	if app.Router != nil {
		p.registerRoutes(app)
	}
	return nil
}

func (p *RoleProvider) registerRoutes(app *kodia.App) {
	roleHandler := kodia.MustResolve[*handlers.RoleHandler](app, "role_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	api := app.Router.Group("/api")
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(jwtManager))
	{
		// Role management
		roles := admin.Group("/roles")
		roles.Use(middleware.RequireRole("admin"))
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.GetRoles)
			roles.DELETE("/:id", roleHandler.DeleteRole)
		}

		// User role assignment
		users := admin.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.POST("/:user_id/roles", roleHandler.AssignRole)
			users.DELETE("/:user_id/roles/:role", roleHandler.RevokeRole)
			users.GET("/:user_id/roles", roleHandler.GetUserRoles)
		}
	}
}
