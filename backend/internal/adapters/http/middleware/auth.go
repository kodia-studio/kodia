// Package middleware contains Gin middleware for Kodia Framework.
package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/response"
)

const (
	userIDKey          = "user_id"
	userEmailKey       = "user_email"
	userRoleKey        = "user_role"
	userPermissionsKey = "user_permissions"
)

// Auth validates the Bearer JWT access token and injects claims into the context.
// Place this middleware on routes that require authentication.
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Invalid authorization header format. Expected: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Inject claims into context for handlers to use
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.Email)
		c.Set(userRoleKey, claims.Role)
		c.Set(userPermissionsKey, claims.Permissions)
		c.Next()
	}
}

// RequireRole returns a middleware that enforces a specific role.
// The Auth middleware must run before this one.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(userRoleKey)
		if !exists {
			response.Forbidden(c, "Authentication required")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.Forbidden(c, "Invalid role context")
			c.Abort()
			return
		}

		for _, role := range roles {
			if strings.EqualFold(roleStr, role) {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient roles")
		c.Abort()
	}
}

// RequirePermission enforces that a user has AT LEAST ONE of the specified permissions.
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPerms, exists := c.Get(userPermissionsKey)
		if !exists {
			response.Forbidden(c, "Authentication required")
			c.Abort()
			return
		}

		perms, ok := userPerms.([]string)
		if !ok {
			response.Forbidden(c, "Invalid permission context")
			c.Abort()
			return
		}

		// Admins have all permissions
		userRole, _ := c.Get(userRoleKey)
		if r, ok := userRole.(string); ok && strings.EqualFold(r, "admin") {
			c.Next()
			return
		}

		for _, required := range permissions {
			for _, userPerm := range perms {
				if userPerm == "*" || strings.EqualFold(required, userPerm) {
					c.Next()
					return
				}
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

// RequireAllPermissions enforces that a user has ALL of the specified permissions.
func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPerms, exists := c.Get(userPermissionsKey)
		if !exists {
			response.Forbidden(c, "Authentication required")
			c.Abort()
			return
		}

		perms, ok := userPerms.([]string)
		if !ok {
			response.Forbidden(c, "Invalid permission context")
			c.Abort()
			return
		}

		for _, required := range permissions {
			found := false
			for _, userPerm := range perms {
				if strings.EqualFold(required, userPerm) {
					found = true
					break
				}
			}
			if !found {
				response.Forbidden(c, fmt.Sprintf("Missing required permission: %s", required))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GetUserID extracts the authenticated user's ID from the Gin context.
func GetUserID(c *gin.Context) string {
	id, _ := c.Get(userIDKey)
	s, _ := id.(string)
	return s
}

// GetUserRole extracts the authenticated user's role from the Gin context.
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get(userRoleKey)
	s, _ := role.(string)
	return s
}

// GetUserPermissions extracts the authenticated user's permissions from the Gin context.
func GetUserPermissions(c *gin.Context) []string {
	perms, _ := c.Get(userPermissionsKey)
	s, _ := perms.([]string)
	return s
}
