package providers

import (
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

type GraphQLProvider struct{}

func NewGraphQLProvider() *GraphQLProvider {
	return &GraphQLProvider{}
}

func (p *GraphQLProvider) Name() string {
	return "kodia:graphql"
}

func (p *GraphQLProvider) Register(app *kodia.App) error {
	authService := app.MustGet("auth_service").(ports.AuthService)
	userService := app.MustGet("user_service").(ports.UserService)
	
	graphqlHandler := handlers.NewGraphQLHandler(authService, userService, app.Log)
	app.Set("graphql_handler", graphqlHandler)

	return nil
}

func (p *GraphQLProvider) Boot(app *kodia.App) error {
	if app.Router != nil {
		h := app.MustGet("graphql_handler").(*handlers.GraphQLHandler)
		api := app.Router.Group("/api")
		api.POST("/query", middleware.GraphQLContextMiddleware(), h.QueryHandler())
		
		if !app.Config.IsProduction() {
			api.GET("/playground", h.PlaygroundHandler())
		}
	}
	return nil
}
