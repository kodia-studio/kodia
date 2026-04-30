# Response Transformers & API Versioning

A comprehensive guide to building version-aware JSON responses and managing API evolution in Kodia Framework.

**Table of Contents**
- [Overview](#overview)
- [Problem Statement](#problem-statement)
- [Architecture](#architecture)
- [Resource Transformers](#resource-transformers)
- [API Versioning](#api-versioning)
- [Deprecation Warnings](#deprecation-warnings)
- [Implementation Guide](#implementation-guide)
- [Examples](#examples)
- [Testing](#testing)

---

## Overview

The Response Transformers & API Versioning system enables:

1. **Type-Safe Transformations** — Convert domain models to JSON using Go generics
2. **Conditional Field Inclusion** — Show/hide fields based on API version
3. **Version-Aware Routing** — Support `/api/v1/`, `/api/v2/` alongside legacy `/api/` endpoints
4. **Backward Compatibility** — Existing unversioned routes continue to work, defaulting to v1
5. **Deprecation Signals** — RFC 8594 headers guide clients to newer API versions
6. **Field Filtering** — Clients can request specific fields via `?fields=id,name,email`

### Key Benefits

- **No Breaking Changes** — Evolve the API without disrupting existing clients
- **Zero Boilerplate** — Generics + middleware eliminate repetitive mapping code
- **Safe Refactoring** — Version-conditional logic allows gradual migration
- **Client Guidance** — Deprecation headers help clients prepare for sunset dates

---

## Problem Statement

### Before

```
❌ Fixed Response Shape
   GET /api/users/123 returns the same fields for all clients
   ├─ Can't add fields without breaking existing clients
   ├─ Can't hide sensitive fields based on context
   └─ Impossible to deprecate fields gradually

❌ No API Versioning
   All routes at /api/* with no version prefix
   ├─ No way to signal v2 API changes to clients
   ├─ Clients forced to parse response content to detect changes
   └─ Can't maintain multiple versions simultaneously

❌ Manual DTOs & Mappers
   Domain Model → DTO → JSON requires verbose mapping:
   │ type User struct { ID, Name, Email, Role, Roles, Permissions }
   ├─ CreateUserDTO (v1 only)
   ├─ UserResponseDTO (v1 only)
   ├─ UpdateUserDTO (v2+)
   └─ MapToResponse() ... (boilerplate)
```

### After

```
✅ Version-Aware Responses
   GET /api/v1/users/123 → base fields
   GET /api/v2/users/123 → includes roles, permissions
   GET /api/v1/users/123?fields=id,name → filtered response

✅ Clean API Versioning
   URL path: /api/v1/, /api/v2/, /api/v3/
   Or header: API-Version: v2
   Or Accept: application/vnd.kodia.v2+json

✅ Type-Safe Transformers
   var User = resource.TransformFunc[*domain.User](func(u *domain.User, ctx resource.TransformContext) map[string]any {
       data := map[string]any{ /* v1 base fields */ }
       if ctx.Since("v2") { data["roles"] = u.Roles }
       return ctx.Only(data)
   })
```

---

## Architecture

### Core Components

#### 1. Resource Transformers (`pkg/resource/resource.go`)

**Transformer Interface** — Generic contract for model → JSON:

```go
type Transformer[T any] interface {
    Transform(model T, ctx TransformContext) map[string]any
}
```

**TransformFunc Adaptor** — Convert a function to a Transformer:

```go
type TransformFunc[T any] func(model T, ctx TransformContext) map[string]any

func (f TransformFunc[T]) Transform(model T, ctx TransformContext) map[string]any {
    return f(model, ctx)
}
```

**Item & Collection Helpers**:

```go
// Single item
func Item[T any](model T, t Transformer[T], ctx TransformContext) map[string]any

// Multiple items
func Collection[T any](models []T, t Transformer[T], ctx TransformContext) []map[string]any

// Extract context from Gin request
func FromContext(c *gin.Context) TransformContext
```

#### 2. Transform Context (`pkg/resource/transform_context.go`)

Metadata for conditional transformations:

```go
type TransformContext struct {
    Version string   // "v1", "v2", etc. (set by version middleware)
    Fields  []string // Optional field whitelist from ?fields=id,name
}
```

**Methods:**

```go
// true if current version >= compare ("v2" >= "v1" = true)
func (tc TransformContext) Since(v string) bool

// true if current version <= compare ("v1" <= "v2" = true)
func (tc TransformContext) Until(v string) bool

// Apply field filter. Empty Fields = return all.
func (tc TransformContext) Only(data map[string]any) map[string]any
```

#### 3. Version Management (`pkg/version/version.go`)

**Middleware** — Extracts and injects API version:

```go
func Middleware() gin.HandlerFunc
```

**Priority Order:**
1. URL path segment: `/api/v2/users` → "v2"
2. Header: `API-Version: v3` → "v3"
3. Accept header: `application/vnd.kodia.v2+json` → "v2"
4. Default: "v1"

**Version Checking** in handlers:

```go
func Get(c *gin.Context) string              // "v1", "v2", etc.
func Since(c *gin.Context, v string) bool    // Is current >= v?
func Until(c *gin.Context, v string) bool    // Is current <= v?
```

**Deprecation Headers** — RFC 8594:

```go
func Deprecate(sunsetDate, alternative string) gin.HandlerFunc
```

Attaches:
- `Deprecated: true`
- `Sunset: Thu, 31 Dec 2026 00:00:00 GMT` (if sunsetDate provided)
- `Link: </api/v2/users>; rel="successor-version"` (if alternative provided)

---

## Resource Transformers

### Defining a Transformer

```go
package transformers

import (
    "time"
    "github.com/kodia-studio/kodia/internal/core/domain"
    "github.com/kodia-studio/kodia/pkg/resource"
)

// User transforms domain.User → JSON response.
// v1: base fields only
// v2+: includes roles and permissions
var User = resource.TransformFunc[*domain.User](func(u *domain.User, ctx resource.TransformContext) map[string]any {
    if u == nil {
        return map[string]any{}
    }

    // Build complete data map
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

    // v2+: add RBAC fields
    if ctx.Since("v2") {
        data["roles"]       = u.Roles
        data["permissions"] = u.Permissions
    }

    // Apply field filter (?fields=id,name)
    return ctx.Only(data)
})
```

### Composite Transformers

Nest transformers for related models:

```go
var Auth = resource.TransformFunc[*ports.AuthResponse](func(r *ports.AuthResponse, ctx resource.TransformContext) map[string]any {
    if r == nil {
        return map[string]any{}
    }

    data := map[string]any{
        "access_token":  r.AccessToken,
        "refresh_token": r.RefreshToken,
        "token_type":    "Bearer",
        "user":          resource.Item(r.User, User, ctx),  // Nested!
    }

    if r.MFARequired {
        data["mfa_required"] = true
        data["mfa_token"]    = r.MFAToken
    }

    return ctx.Only(data)
})
```

### Conditional Fields

Use `Since()` and `Until()` for version-dependent inclusion:

```go
// Only in v1-v2
if ctx.Until("v2") {
    data["deprecated_field"] = u.OldField
}

// Starting from v2
if ctx.Since("v2") {
    data["new_field"] = u.NewField
}

// Exact version
if ctx.Version == "v3" {
    data["v3_only"] = u.V3Data
}
```

---

## API Versioning

### Dual Route Registration

Register routes on both versioned and legacy paths:

```go
func (p *AuthProvider) registerRoutes(app *kodia.App) {
    authHandler := kodia.MustResolve[*handlers.AuthHandler](app, "auth_handler")
    jwtManager  := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

    api := app.Router.Group("/api")

    // New: versioned routes with middleware
    v1 := api.Group("/v1")
    v1.Use(version.Middleware())
    registerAuthRoutes(v1.Group("/auth"), authHandler, jwtManager)

    // Legacy: unversioned routes (backward compatibility, defaults to v1)
    registerAuthRoutes(api.Group("/auth"), authHandler, jwtManager)
}

func registerAuthRoutes(g *gin.RouterGroup, h *handlers.AuthHandler, jwtManager *jwt.Manager) {
    g.POST("/register", h.Register)
    g.POST("/login", h.Login)
    g.POST("/refresh", h.RefreshToken)
    g.GET("/verify-email", h.VerifyEmail)
    // ... more routes
}
```

This creates:
- `POST /api/v1/auth/login` ← new, versioned
- `POST /api/auth/login` ← legacy, backward compat

### Version Detection Methods

**1. URL Path** (highest priority)
```
GET /api/v2/users/123 → version = "v2"
GET /api/v1/users/123 → version = "v1"
```

**2. Header**
```
GET /api/users/123 \
  -H "API-Version: v2" → version = "v2"
```

**3. Accept Header**
```
GET /api/users/123 \
  -H "Accept: application/vnd.kodia.v2+json" → version = "v2"
```

**4. Default (if none above)**
```
GET /api/users/123 → version = "v1"
```

---

## Deprecation Warnings

### Marking Endpoints Deprecated

```go
func (p *UserProvider) registerRoutes(app *kodia.App) {
    api := app.Router.Group("/api")

    // Legacy v1 endpoint, sunset 2026-12-31
    v1 := api.Group("/v1")
    v1.Use(version.Middleware())
    v1.GET("/users/:id",
        version.Deprecate("2026-12-31", "/api/v2/users/:id"),
        h.GetByID,
    )

    // Current version
    v2 := api.Group("/v2")
    v2.Use(version.Middleware())
    v2.GET("/users/:id", h.GetByID)
}
```

### RFC 8594 Headers

The client receives:

```
HTTP/1.1 200 OK
Deprecated: true
Sunset: Thu, 31 Dec 2026 00:00:00 GMT
Link: </api/v2/users/123>; rel="successor-version"
Content-Type: application/json

{ "id": "123", "name": "John" }
```

**Client Interpretation:**
- `Deprecated: true` — This endpoint is deprecated
- `Sunset: ...` — Stop using this endpoint after this date
- `Link: ...` — Use this URL instead

---

## Implementation Guide

### Step 1: Create a Transformer

File: `internal/adapters/http/transformers/my_transformer.go`

```go
package transformers

import (
    "github.com/kodia-studio/kodia/internal/core/domain"
    "github.com/kodia-studio/kodia/pkg/resource"
)

var MyModel = resource.TransformFunc[*domain.MyModel](func(m *domain.MyModel, ctx resource.TransformContext) map[string]any {
    if m == nil {
        return map[string]any{}
    }

    data := map[string]any{
        "id":   m.ID,
        "name": m.Name,
    }

    if ctx.Since("v2") {
        data["extra"] = m.Extra
    }

    return ctx.Only(data)
})
```

### Step 2: Use in Handler

```go
func (h *MyHandler) GetMyModel(c *gin.Context) {
    // ... fetch model ...

    ctx := resource.FromContext(c)
    response.OK(c, "Success", resource.Item(model, transformers.MyModel, ctx))
}
```

### Step 3: Register Versioned Routes

```go
func (p *MyProvider) registerRoutes(app *kodia.App) {
    handler := kodia.MustResolve[*handlers.MyHandler](app, "my_handler")

    api := app.Router.Group("/api")

    v1 := api.Group("/v1")
    v1.Use(version.Middleware())
    registerMyRoutes(v1.Group("/my"), handler)

    registerMyRoutes(api.Group("/my"), handler)
}

func registerMyRoutes(g *gin.RouterGroup, h *handlers.MyHandler) {
    g.GET("/:id", h.GetMyModel)
    g.GET("", h.ListMyModels)
}
```

### Step 4: Test

```go
func TestMyModelTransformer(t *testing.T) {
    transformer := transformers.MyModel

    // v1 response
    ctx := resource.TransformContext{Version: "v1"}
    result := resource.Item(model, transformer, ctx)
    if _, ok := result["extra"]; ok {
        t.Error("v1 should not have 'extra' field")
    }

    // v2 response
    ctx = resource.TransformContext{Version: "v2"}
    result = resource.Item(model, transformer, ctx)
    if _, ok := result["extra"]; !ok {
        t.Error("v2 should have 'extra' field")
    }
}
```

---

## Examples

### Example 1: Simple Versioning

**GET /api/v1/users/123**
```json
{
  "id": "123",
  "name": "Alice",
  "email": "alice@example.com"
}
```

**GET /api/v2/users/123**
```json
{
  "id": "123",
  "name": "Alice",
  "email": "alice@example.com",
  "roles": ["admin", "editor"],
  "permissions": ["users:read", "users:write"]
}
```

### Example 2: Field Filtering

**GET /api/v1/users/123?fields=id,name**
```json
{
  "id": "123",
  "name": "Alice"
}
```

### Example 3: Composite Response

**GET /api/v1/auth/login**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "user": {
    "id": "123",
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

### Example 4: Deprecation Migration

**GET /api/v1/old-endpoint**
```
HTTP/1.1 200 OK
Deprecated: true
Sunset: Sun, 30 Jun 2026 00:00:00 GMT
Link: </api/v2/new-endpoint>; rel="successor-version"

{ "data": ... }
```

Client receives clear signal to migrate before June 30, 2026.

---

## Testing

### Unit Testing Transformers

```go
func TestUserTransformer(t *testing.T) {
    tests := []struct {
        name     string
        version  string
        expected bool // should have "roles"?
    }{
        {"v1 excludes roles", "v1", false},
        {"v2 includes roles", "v2", true},
        {"v3 includes roles", "v3", true},
    }

    user := &domain.User{
        ID:    "123",
        Name:  "Alice",
        Roles: []string{"admin"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := resource.TransformContext{Version: tt.version}
            result := resource.Item(user, transformers.User, ctx)

            if _, ok := result["roles"]; ok != tt.expected {
                t.Errorf("roles presence: got %v, want %v", ok, tt.expected)
            }
        })
    }
}
```

### Integration Testing with Routes

```go
func TestUserEndpointVersioning(t *testing.T) {
    router := setupRouter()

    // Test v1
    resp := performRequest(router, "GET", "/api/v1/users/123")
    if status := resp.Code; status != 200 {
        t.Errorf("v1: got status %d", status)
    }
    var data map[string]any
    json.Unmarshal(resp.Body.Bytes(), &data)
    if _, hasRoles := data["data"].(map[string]any)["roles"]; hasRoles {
        t.Error("v1 should not have roles field")
    }

    // Test v2
    resp = performRequest(router, "GET", "/api/v2/users/123")
    json.Unmarshal(resp.Body.Bytes(), &data)
    if _, hasRoles := data["data"].(map[string]any)["roles"]; !hasRoles {
        t.Error("v2 should have roles field")
    }
}
```

---

## FAQ

**Q: Can I use both transformers and DTOs?**
A: Yes. Migrate gradually — some endpoints use transformers, others use DTOs. No requirement to convert everything at once.

**Q: What happens if a client requests /api/users/123 (no version)?**
A: The unversioned route acts as v1 (default). The version middleware sets "v1" in context.

**Q: Can I deprecate a field instead of an endpoint?**
A: Yes, using `ctx.Until("v1")` to exclude it from v2+:
```go
if ctx.Until("v1") {
    data["old_field"] = m.OldField  // Only in v1
}
```

**Q: How do I handle breaking changes?**
A: Version-conditional logic:
```go
if ctx.Since("v2") {
    data["new_field"] = m.NewField   // v2+
} else {
    data["old_field"] = m.OldField   // v1 only
}
```

**Q: Can multiple clients use different API versions simultaneously?**
A: Yes. Clients explicitly choose via URL path, header, or Accept. The framework supports all versions at once.

**Q: How do I test field filtering (?fields=)?**
A: The TransformContext.Fields is populated from query params:
```go
ctx := resource.FromContext(c)  // Reads ?fields=id,name from request
```

---

## Best Practices

1. **Always default to v1** — New clients start here
2. **Plan deprecation early** — Signal sunset dates at least 3 months in advance
3. **Version-condition at point of change** — Keep logic close to the field
4. **Test both versions** — Unit tests should verify version differences
5. **Document field changes** — Note in API docs when fields were added/removed
6. **Keep route structure consistent** — Same paths across versions
7. **Use field filtering sparingly** — ?fields= is for optimization, not permissions

---

## Troubleshooting

**Issue: Version always resolves to "v1"**
- Check that version.Middleware() is applied to the router group
- Verify the URL path contains /v2/, /v3/, etc.
- Try explicit header: `-H "API-Version: v2"`

**Issue: Fields missing from response**
- Ensure transformer checks `ctx.Since()` / `ctx.Until()` correctly
- Verify `ctx.Only()` is called (if Fields whitelist is set)
- Check that the field is actually in the data map

**Issue: Deprecation headers not sent**
- Ensure `version.Deprecate()` middleware is before handler
- Verify sunsetDate format is "YYYY-MM-DD"
- Check that both Deprecated and Sunset headers are present

**Issue: Nested transformer returns nil**
- Ensure nested transformer handles nil inputs gracefully
- Return empty map `map[string]any{}` on nil, not nil itself

