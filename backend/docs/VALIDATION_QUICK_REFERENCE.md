# Validation Layer — Quick Reference Guide

## TL;DR

```go
// 1. Define DTO with validation rules
type RegisterRequest struct {
    Name     string `json:"name"     validate:"required,min=2,max=100"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,strong_password"`
}

// 2. Use in handler (3 lines!)
func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    if !validation.BindAndValidate(c, h.validate, &req) {
        return  // Error already sent
    }
    // req is valid, continue...
}
```

## Validation Rules Cheat Sheet

### Kodia Custom Rules

```
strong_password   # uppercase + lowercase + number + symbol
phone             # valid phone number (E.164)
alpha_space       # only letters and spaces
no_html           # no < or > characters
```

### Go-Playground Rules (Most Common)

```
required          # field required
email             # valid email
min=8             # min 8 chars/value
max=100           # max 100 chars/value
len=6             # exactly 6 chars
url               # valid URL
uuid4             # valid UUID
oneof=a b c       # one of these values
gt=0              # greater than 0
gte=0             # greater than or equal
lte=130           # less than or equal
omitempty         # skip if not provided
```

## Common DTO Patterns

### Password Field
```go
Password string `json:"password" validate:"required,min=8,max=72,strong_password"`
```

### Email Field
```go
Email string `json:"email" validate:"required,email"`
```

### Name Field
```go
Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
```

### Optional Field
```go
AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
```

### User-Generated Content (XSS Prevention)
```go
Bio string `json:"bio" validate:"required,max=500,no_html"`
```

### Multiple Rules
```go
Title string `json:"title" validate:"required,min=5,max=200,no_html"`
```

## Handler Pattern

```go
func (h *MyHandler) HandleRequest(c *gin.Context) {
    // 1. Declare request struct
    var req dto.MyRequest
    
    // 2. Bind + Validate (all in one!)
    if !validation.BindAndValidate(c, h.validate, &req) {
        return  // Error response sent automatically
    }
    
    // 3. Use validated request
    result, err := h.service.DoSomething(c.Request.Context(), req)
    
    // 4. Handle service errors
    if err != nil {
        // Handle error...
        return
    }
    
    // 5. Return success
    response.OK(c, "Success", result)
}
```

## Error Responses

### Validation Error (422)
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

### Binding Error (400)
```json
{
    "success": false,
    "message": "Invalid request body"
}
```

## Key Points

| Concept | Details |
|---------|---------|
| **Format Validation** | Validation layer (regex, format, size) |
| **Business Logic Validation** | Service layer (email taken, permission denied) |
| **Error Response** | Automatic (422 with per-field messages) |
| **Handler Boilerplate** | Reduced from ~10 lines to 3 lines |
| **Custom Rules** | strong_password, phone, alpha_space, no_html |
| **Optional Fields** | Use `omitempty` to skip validation |

## Do's and Don'ts

### ✓ Do

- Use `validation.BindAndValidate()` in all handlers
- Define validation rules in DTOs via struct tags
- Use `strong_password` for all password fields
- Use `no_html` for user-generated content
- Use `alpha_space` for name fields
- Use service layer for business logic validation

### ✗ Don't

- Manually call `c.ShouldBindJSON()` and `h.validate.Struct()` separately
- Put validation rules in handler code
- Use weak password validation (min length only)
- Allow HTML in user input without validation
- Mix format validation (validation layer) with business logic (service layer)

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "validation not found" | Import: `github.com/kodia-studio/kodia/pkg/validation` |
| Validation not triggering | Use `validate:` tags, not `binding:` tags |
| Custom rules not working | Make sure to use `validation.New()` |
| Error response not sent | `BindAndValidate` sends response automatically on error |
| Field name wrong in errors | Check JSON tag in DTO struct |

## Full Example

```go
package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/kodia-studio/kodia/internal/adapters/http/dto"
    "github.com/kodia-studio/kodia/pkg/response"
    "github.com/kodia-studio/kodia/pkg/validation"
)

type UserHandler struct {
    userService UserService
    validate    *validation.Validator
}

func (h *UserHandler) Create(c *gin.Context) {
    // 1. Validate
    var req dto.CreateUserRequest
    if !validation.BindAndValidate(c, h.validate, &req) {
        return
    }

    // 2. Create user
    user, err := h.userService.Create(c.Request.Context(), req)
    if err != nil {
        response.InternalServerError(c, "Failed to create user")
        return
    }

    // 3. Return success
    response.Created(c, "User created", user)
}
```

## Next Steps

- Read [Validation Layer Documentation](VALIDATION_LAYER.md) for detailed guide
- Check [kodia-web Integration Guide](../kodia-web/backend/docs/VALIDATION_LAYER.md)
- View [Full Validation Examples](../internal/adapters/http/dto/dto.go)
