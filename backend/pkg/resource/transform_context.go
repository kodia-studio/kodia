package resource

import (
	"strconv"
	"strings"
)

// TransformContext carries per-request context into a Transformer.
type TransformContext struct {
	Version string   // "v1", "v2", etc. — from version.Middleware
	Fields  []string // Optional field whitelist from ?fields=id,name,email
}

// Since returns true if the current API version >= v.
// Example: ctx.Since("v2") is true when client uses v2 or higher.
func (tc TransformContext) Since(v string) bool {
	return versionNum(tc.Version) >= versionNum(v)
}

// Until returns true if the current API version <= v.
// Example: ctx.Until("v1") is true only for v1 clients.
func (tc TransformContext) Until(v string) bool {
	return versionNum(tc.Version) <= versionNum(v)
}

// Only applies the field whitelist. If Fields is empty, returns data unchanged.
// Usage in transformer: return ctx.Only(map[string]any{ "id": u.ID, ... })
func (tc TransformContext) Only(data map[string]any) map[string]any {
	if len(tc.Fields) == 0 {
		return data
	}
	result := make(map[string]any, len(tc.Fields))
	for _, f := range tc.Fields {
		f = strings.TrimSpace(f)
		if val, ok := data[f]; ok {
			result[f] = val
		}
	}
	return result
}

// versionNum extracts the numeric part from "v1", "v2", etc.
func versionNum(v string) int {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.TrimPrefix(v, "v")
	n, _ := strconv.Atoi(v)
	return n
}
