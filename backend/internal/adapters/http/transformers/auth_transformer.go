package transformers

import (
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/resource"
)

// Auth transforms a ports.AuthResponse to a JSON map.
var Auth = resource.TransformFunc[*ports.AuthResponse](func(r *ports.AuthResponse, ctx resource.TransformContext) map[string]any {
	if r == nil {
		return map[string]any{}
	}

	data := map[string]any{
		"access_token":  r.AccessToken,
		"refresh_token": r.RefreshToken,
		"token_type":    "Bearer",
		"user":          resource.Item(r.User, User, ctx),
	}

	if r.MFARequired {
		data["mfa_required"] = true
		data["mfa_token"]    = r.MFAToken
	}

	return ctx.Only(data)
})
