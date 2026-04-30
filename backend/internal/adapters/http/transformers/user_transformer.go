package transformers

import (
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/resource"
)

// User transforms a domain.User to a JSON map.
// v1: base fields (id, name, email, role, is_active, etc.)
// v2+: adds "roles" (RBAC roles slice) and "permissions" slice
var User = resource.TransformFunc[*domain.User](func(u *domain.User, ctx resource.TransformContext) map[string]any {
	if u == nil {
		return map[string]any{}
	}

	data := map[string]any{
		"id":                 u.ID,
		"name":               u.Name,
		"email":              u.Email,
		"role":               string(u.Role),
		"is_active":          u.IsActive,
		"is_verified":        u.IsVerified,
		"two_factor_enabled": u.TwoFactorEnabled,
		"avatar_url":         u.AvatarURL,
		"created_at":         u.CreatedAt.Format(time.RFC3339),
		"updated_at":         u.UpdatedAt.Format(time.RFC3339),
	}

	// v2+: RBAC roles and permissions
	if ctx.Since("v2") {
		data["roles"]       = u.Roles
		data["permissions"] = u.Permissions
	}

	return ctx.Only(data)
})
