package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/tenancy"
)

/**
 * TenantMiddleware identifies the current tenant.
 */
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Identification via Custom Header
		tenantID := c.GetHeader("X-Tenant-ID")

		// 2. Identification via Query (optional convenience)
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}

		if tenantID == "" {
			// In a real SaaS, we might block access if no tenant is provided
			// For general framework usage, we might allow it (default tenant)
		}

		// Inject into context
		c.Set(tenancy.ContextTenantKey, tenantID)

		// 3. Super Admin Detection (Bypass logic)
		// Usually set by AuthMiddleware after verifying JWT
		// Example: user := c.Get("user").(*domain.User)
		// For this implementation, we check for a specific header or claim simulation
		if c.GetHeader("X-Super-Admin") == "true" {
			c.Set(tenancy.ContextIsAdminKey, true)
		}

		c.Next()
	}
}
