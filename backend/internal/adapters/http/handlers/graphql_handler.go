package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/graph"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// GraphQLHandler wraps gqlgen server for Gin integration.
type GraphQLHandler struct {
	authService ports.AuthService
	userService ports.UserService
	log         *zap.Logger
}

// NewGraphQLHandler creates a new GraphQLHandler.
func NewGraphQLHandler(authService ports.AuthService, userService ports.UserService, log *zap.Logger) *GraphQLHandler {
	return &GraphQLHandler{
		authService: authService,
		userService: userService,
		log:         log,
	}
}

// QueryHandler defines the main GraphQL query endpoint.
func (h *GraphQLHandler) QueryHandler() gin.HandlerFunc {
	// Root resolver with dependencies
	res := &graph.Resolver{
		AuthService: h.authService,
		UserService: h.userService,
		Log:         h.log,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: res}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler defines the GraphQL Playground UI endpoint.
func (h *GraphQLHandler) PlaygroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
