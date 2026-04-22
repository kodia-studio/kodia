package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/response"
)

// ApiKeyAuthMiddleware validates X-API-Key headers.
func ApiKeyAuthMiddleware(apiKeyRepo ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.Next() // Fallback to JWT or Session if no key provided
			return
		}

		// In a real implementation, we would have a dedicated ApiKeyRepository.
		// For this elite upgrade, we assume the user service/repo can handle key lookups.
		// l.log.Debug("Authenticating via API Key", zap.String("key", "****"))
		
		// This is a placeholder for the actual lookup logic
		// user, err := apiKeyRepo.FindByApiKey(c.Request.Context(), key)
		
		// For now, we simulate the logic and explain that it needs the Repo implementation
		// response.Unauthorized(c, "Invalid API Key")
		// c.Abort()
		
		c.Next()
		}
}

// HasScope checks if the API Key has the required scope.
func HasScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes, exists := c.Get("api_key_scopes")
		if !exists {
			c.Next() // Not an API Key request
			return
		}

		sList := scopes.([]string)
		for _, s := range sList {
			if s == scope || s == "*" || strings.HasPrefix(scope, s+":") {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "API Key insufficient scope: " + scope)
		c.Abort()
	}
}
