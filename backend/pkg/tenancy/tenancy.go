package tenancy

import "github.com/gin-gonic/gin"

const (
	ContextTenantKey   = "kodia:tenant_id"
	ContextIsAdminKey  = "kodia:is_super_admin"
)

/**
 * Tenantable is an interface that models can implement to tell Kodia
 * that they belong to a specific tenant.
 */
type Tenantable interface {
	GetTenantID() string
}

/**
 * GetTenantID helper to retrieve the tenant ID from context.
 */
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(ContextTenantKey); exists {
		return tenantID.(string)
	}
	return ""
}

/**
 * IsSuperAdmin helper to check if the current user bypasses tenant isolation.
 */
func IsSuperAdmin(c *gin.Context) bool {
	if isAdmin, exists := c.Get(ContextIsAdminKey); exists {
		return isAdmin.(bool)
	}
	return false
}
