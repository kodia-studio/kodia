// Package resource provides generic, version-aware resource transformers for HTTP responses.
// Transformers convert domain models to JSON-serializable maps, with support for
// API versioning, field filtering, and conditional field inclusion.
package resource

import "github.com/gin-gonic/gin"

// Transformer converts a model of type T to a JSON-serializable map.
// The TransformContext carries the API version and requested fields.
type Transformer[T any] interface {
	Transform(model T, ctx TransformContext) map[string]any
}

// TransformFunc is a function satisfying Transformer[T].
type TransformFunc[T any] func(model T, ctx TransformContext) map[string]any

func (f TransformFunc[T]) Transform(model T, ctx TransformContext) map[string]any {
	return f(model, ctx)
}

// Item transforms a single model using the given transformer.
func Item[T any](model T, t Transformer[T], ctx TransformContext) map[string]any {
	return t.Transform(model, ctx)
}

// Collection transforms a slice of models.
func Collection[T any](models []T, t Transformer[T], ctx TransformContext) []map[string]any {
	result := make([]map[string]any, len(models))
	for i, m := range models {
		result[i] = t.Transform(m, ctx)
	}
	return result
}

// FromContext extracts a TransformContext from a gin.Context.
// Reads the API version injected by version.Middleware and optional ?fields= query param.
func FromContext(c *gin.Context) TransformContext {
	v, _ := c.Get("api_version")
	ver, _ := v.(string)
	if ver == "" {
		ver = "v1"
	}
	return TransformContext{
		Version: ver,
		Fields:  c.QueryArray("fields"),
	}
}
