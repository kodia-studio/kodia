package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// A private key for context that only this package can access.
type key string

const (
	userCtxKey key = "user"
)

// UserData holds information injected into GraphQL context
type UserData struct {
	ID    string
	Email string
	Role  string
}

// GraphQLContextMiddleware injects user data from Gin context into standard context
func GraphQLContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(userIDKey)
		userEmail, _ := c.Get(userEmailKey)
		userRole, _ := c.Get(userRoleKey)

		if userID != "" {
			user := &UserData{
				ID:    userID.(string),
				Email: userEmail.(string),
				Role:  userRole.(string),
			}
			// Put it in context
			ctx := context.WithValue(c.Request.Context(), userCtxKey, user)
			// and return it in the request
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserData {
	raw, _ := ctx.Value(userCtxKey).(*UserData)
	return raw
}
