# Request Validation Layer

The Kodia validation layer provides a centralized, reusable way to validate incoming HTTP requests. It wraps `go-playground/validator` with Kodia-specific features and eliminates boilerplate code from handlers.

## Overview

Instead of repeating bind + validate + format errors in every handler, use the validation layer:

```go
// ✓ Clean: 3 lines of validation code
var req dto.RegisterRequest
if !validation.BindAndValidate(c, h.validate, &req) {
    return  // Error response already sent by BindAndValidate
}
// Handler logic continues...
```

## Core Components

### 1. Validator

The `Validator` wraps `go-playground/validator` with automatic custom rule registration and JSON field name support.

```go
// Create a validator (usually in providers)
validator := validation.New()

// Use in handler
if !validation.BindAndValidate(c, validator, &req) {
    return
}
```

**Features:**
- All custom rules pre-registered
- JSON field names in error messages ("email" not "Email")
- Backward compatible with `validator.Struct()`

### 2. BindAndValidate Helper

Convenience function that combines binding and validation in one call:

```go
func BindAndValidate(c *gin.Context, vl *Validator, req any) bool
```

**Behavior:**
- Binds JSON body to `req` via `c.ShouldBindJSON()`
- Validates using the validator
- Returns `true` if successful
- Returns `false` if binding or validation fails (error response already sent)

**Response on failure:**
- Binding error → 400 Bad Request
- Validation error → 422 Unprocessable Entity with per-field errors

### 3. Generic Middleware (Optional)

For opt-in route-level validation:

```go
// In router setup
router.POST("/register", validation.Middleware[dto.RegisterRequest](validator), handler.Register)

// In handler - get already-validated request
func (h *Handler) Register(c *gin.Context) {
    req := validation.Get[dto.RegisterRequest](c)
    // req is already validated, no need to check errors
}
```

## Validation Rules

### Built-in Rules (from go-playground/validator)

| Tag | Description | Example |
|---|---|---|
| `required` | Field must not be empty | `validate:"required"` |
| `email` | Valid email address | `validate:"email"` |
| `min=N` | Minimum length/value | `validate:"min=8"` |
| `max=N` | Maximum length/value | `validate:"max=100"` |
| `len=N` | Exact length | `validate:"len=6"` |
| `url` | Valid URL | `validate:"url"` |
| `uuid4` | Valid UUID v4 | `validate:"uuid4"` |
| `oneof=a b c` | One of specified values | `validate:"oneof=admin user"` |
| `gte=N` | Greater than or equal | `validate:"gte=0"` |
| `lte=N` | Less than or equal | `validate:"lte=130"` |

### Custom Rules (Kodia-specific)

| Tag | Description | Example |
|---|---|---|
| `strong_password` | Uppercase, lowercase, digit, symbol | `validate:"strong_password"` |
| `phone` | Valid phone number (E.164 format) | `validate:"phone"` |
| `alpha_space` | Letters and spaces only | `validate:"alpha_space"` |
| `no_html` | No HTML tags (`<` or `>`) | `validate:"no_html"` |

## Usage Examples

### Basic Validation in Handler

```go
package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/kodia-studio/kodia/internal/adapters/http/dto"
    "github.com/kodia-studio/kodia/pkg/validation"
)

type AuthHandler struct {
    validate *validation.Validator
    // ... other fields
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if !validation.BindAndValidate(c, h.validate, &req) {
        return  // Error already sent
    }

    // req is guaranteed to be valid here
    user, err := h.authService.Register(c.Request.Context(), req)
    // ... rest of handler
}
```

### DTO with Multiple Rules

```go
package dto

type CreateProductRequest struct {
    Name        string  `json:"name"        validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"required,min=10,max=1000,no_html"`
    Price       float64 `json:"price"       validate:"required,gt=0"`
    Phone       string  `json:"phone"       validate:"omitempty,phone"`
    Email       string  `json:"email"       validate:"required,email"`
}
```

### Validation Error Response

When validation fails, the client receives:

```json
{
    "success": false,
    "message": "Validation failed",
    "errors": {
        "email": ["email must be a valid email address"],
        "password": [
            "password is too short (min 8 chars)",
            "password must contain uppercase, lowercase, number, and symbol"
        ]
    }
}
```

## Integration in Providers

When creating handlers in your ServiceProvider, pass a validator:

```go
package providers

import "github.com/kodia-studio/kodia/pkg/validation"

func (p *AuthProvider) Register(app *kodia.App) error {
    // Create validator with custom rules
    validate := validation.New()

    // Pass to handler
    authHandler := handlers.NewAuthHandler(authService, validate, app.Log)
    app.Set("auth_handler", authHandler)

    return nil
}
```

## Error Messages

The validation layer provides user-friendly error messages for all rules:

```
required        → "field is required"
email           → "field must be a valid email address"
min             → "field is too short (min 8 chars)"
max             → "field is too long (max 100 chars)"
len             → "field must be exactly 6 characters"
url             → "field must be a valid URL"
uuid4           → "field must be a valid UUID"
strong_password → "field must contain uppercase, lowercase, number, and symbol"
phone           → "field must be a valid phone number"
alpha_space     → "field must contain only letters and spaces"
no_html         → "field must not contain HTML tags"
```

## Advanced Usage

### Custom Error Messages

While the default messages are good, you can customize them by wrapping `validation.FormatErrors()`:

```go
func formatCustomErrors(err error) map[string][]string {
    errs := validation.FormatErrors(err)
    
    // Customize specific errors
    if msgs, ok := errs["password"]; ok {
        errs["password"] = []string{"Password must be strong (8+ chars, mixed case, numbers, symbols)"}
    }
    
    return errs
}
```

### Conditional Validation

Use `omitempty` to allow optional fields:

```go
type UpdateUserRequest struct {
    Name      *string `json:"name"       validate:"omitempty,min=2,max=100"`
    AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
}
```

The field is only validated if it's provided.

### Multiple Validators

You can register multiple validators if needed (advanced):

```go
v := validation.New()  // Has custom rules
v2 := validation.New() // Separate instance, also has custom rules
```

But in practice, a single shared validator per handler is sufficient.

## Best Practices

### 1. Always Use BindAndValidate

❌ **Avoid:**
```go
if err := c.ShouldBindJSON(&req); err != nil { ... }
if err := h.validate.Struct(req); err != nil { ... }
```

✅ **Do:**
```go
if !validation.BindAndValidate(c, h.validate, &req) { return }
```

### 2. Put Validation Rules in DTOs

Validation rules live in your DTO struct tags, not in handlers or services:

```go
// ✓ Good: Rules in DTO
type UserRequest struct {
    Email string `json:"email" validate:"required,email"`
}

// ✗ Avoid: Rules scattered in handler
var email string
if email == "" { /* error */ }
if !strings.Contains(email, "@") { /* error */ }
```

### 3. Business Logic in Services

The validator checks **format** (syntax validation). Use the **service layer** for **business logic** (semantic validation):

```go
// Handler: format validation
var req dto.RegisterRequest
if !validation.BindAndValidate(c, h.validate, &req) {
    return
}

// Service: business logic validation
user, err := h.authService.Register(ctx, req)
if err == domain.ErrEmailAlreadyExists {
    // Email format was valid, but it's already taken (business rule)
    response.Conflict(c, "Email already registered")
    return
}
```

### 4. Use Strong Password Validation for Security-Sensitive Fields

```go
type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password"     validate:"required,strong_password"`
}
```

This ensures all password changes enforce strong password policy.

## Testing Validation

When testing handlers, you can test validation separately:

```go
func TestRegisterValidation(t *testing.T) {
    v := validation.New()
    
    req := dto.RegisterRequest{
        Email:    "invalid-email",
        Password: "weak",
    }
    
    errors := validation.FormatErrors(v.Struct(req))
    
    assert.Contains(t, errors["email"], "must be a valid email address")
    assert.Contains(t, errors["password"], "too short")
}
```

## Migration from Old Validation Approach

If you have old handlers with boilerplate validation code:

**Before:**
```go
var req dto.SomeRequest
if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, "Invalid request body", nil)
    return
}
if err := h.validate.Struct(req); err != nil {
    response.BadRequest(c, "Validation failed", formatValidationErrors(err))
    return
}
```

**After:**
```go
var req dto.SomeRequest
if !validation.BindAndValidate(c, h.validate, &req) {
    return
}
```

Just replace those ~10 lines with 3 lines. The rest of the handler stays the same.

## See Also

- [Request Handling Guide](ARCHITECTURE.md#http-layer)
- [Error Handling Best Practices](JWT_SECURITY.md#error-responses)
- [Service Layer Design](ARCHITECTURE.md#service-layer)
