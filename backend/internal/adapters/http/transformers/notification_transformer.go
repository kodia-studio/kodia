package transformers

import (
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/pkg/resource"
)

// Notification transforms a domain.Notification to a JSON map.
// v1: base fields (id, type, title, message, is_read, created_at)
// v2+: includes "data" field (extra payload)
var Notification = resource.TransformFunc[*domain.Notification](func(n *domain.Notification, ctx resource.TransformContext) map[string]any {
	if n == nil {
		return map[string]any{}
	}

	data := map[string]any{
		"id":         n.ID,
		"type":       string(n.Type),
		"title":      n.Title,
		"message":    n.Message,
		"is_read":    n.IsRead,
		"created_at": n.CreatedAt.Format(time.RFC3339),
	}

	// v2+: include extra data payload
	if ctx.Since("v2") && len(n.Data) > 0 {
		data["data"] = n.Data
	}

	return ctx.Only(data)
})
