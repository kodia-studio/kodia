# Data Validation & Sanitization Framework

## Overview

The Data Validation & Sanitization Framework is a comprehensive system for validating and cleaning user input in the Kodia Framework. It combines declarative struct tags, custom validation rules, asynchronous context-aware validation, and localized error messages—all with **zero new dependencies**.

### Key Features

- **Declarative Validation** — Validation rules via `validate:` struct tags
- **Input Sanitization** — Automatic cleaning via `sanitize:` struct tags with 10+ filters
- **Context-Aware Validation** — Async rules (database lookups) with `ctx:` prefix
- **Localized Error Messages** — English and Indonesian (extensible)
- **Type-Safe Binding** — Unified pipeline with consistent HTTP semantics

### Quick Example

```go
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100,alpha_space" sanitize:"trim"`
    Email    string `json:"email" validate:"required,email,ctx:unique_email" sanitize:"trim,lowercase"`
    Password string `json:"password" validate:"required,min=8,max=72,strong_password"`
}

// In handler
if err := binder.Bind(c, &req); err != nil {
    return  // 400 (JSON error) or 422 (validation error) already sent
}
// req.Name is trimmed, req.Email is lowercase, all validated
```

---

## Architecture

### Component Stack

```
HTTP Request (JSON)
    ↓
[JSON Binding] — c.ShouldBindJSON()
    ↓
[Sanitization] — Apply sanitize: filters (trim, lowercase, etc.)
    ↓
[Sync Validation] — Apply validate: rules (required, email, custom rules)
    ↓
[Async Validation] — Apply ctx: rules (database checks)
    ↓
[Error Formatting] — Localized error messages
    ↓
Response (422 Unprocessable Entity with per-field errors)
```

### Package Structure

```
pkg/validation/
├── sanitizer.go        # Sanitize() + filter pipeline
├── validator.go        # Validator + context-aware rules
├── rules.go            # Custom validation rules (strong_password, phone, etc.)
├── middleware.go       # BindAndValidate(), BindAndValidateCtx()
└── locale.go           # Localized error messages (EN, ID)

pkg/binder/
└── binder.go           # Bind() and BindCtx() shortcuts
```

---

## Sanitization Layer

### Overview

Sanitization cleans user input **before** validation. It's applied via `sanitize:` struct tags on string fields.

### Available Filters

| Filter | Description | Example |
|--------|-------------|---------|
| `trim` | Remove leading/trailing whitespace | `"  hello  "` → `"hello"` |
| `lowercase` | Convert to lowercase | `"HELLO"` → `"hello"` |
| `uppercase` | Convert to uppercase | `"hello"` → `"HELLO"` |
| `escape_html` | Escape HTML entities | `"<p>"` → `"&lt;p&gt;"` |
| `strip_html` | Remove all HTML tags | `"<p>text</p>"` → `"text"` |
| `slug` | URL-safe slug format | `"Hello World"` → `"hello-world"` |
| `no_spaces` | Remove all whitespace | `"hello world"` → `"helloworld"` |
| `alphanumeric` | Keep only letters and digits | `"hello123@"` → `"hello123"` |
| `normalize` | Unicode NFC normalization | `"café"` → `"café"` (normalized) |
| `truncate:N` | Truncate to N characters | `"hello"` with `truncate:3` → `"hel"` |

### Usage Examples

```go
type CreateBlogRequest struct {
    // Trim spaces from title
    Title string `json:"title" sanitize:"trim"`

    // Lowercase email, trim spaces
    Email string `json:"email" sanitize:"trim,lowercase"`

    // Create URL-safe slug
    Slug string `json:"slug" sanitize:"slug"`

    // Remove HTML, keep only text
    Description string `json:"description" sanitize:"strip_html,trim"`

    // Remove special chars from username
    Username string `json:"username" sanitize:"lowercase,alphanumeric"`
}
```

### Chaining Filters

Filters are applied left-to-right:

```go
type UserRequest struct {
    // 1. trim → "  Hello World  " → "Hello World"
    // 2. lowercase → "hello world"
    // 3. slug → "hello-world"
    Slug string `json:"slug" sanitize:"trim,lowercase,slug"`
}
```

### API

```go
// Sanitize applies sanitize: tags to string fields in a struct
// v must be a pointer to a struct; non-string fields are skipped
func Sanitize(v any) error {
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
        return nil  // Silently skip non-struct pointers
    }
    // ... applies filters
}
```

---

## Validation Layer

### Sync Validation

Sync validation rules are checked immediately and don't require context. They use the `validate:` tag.

#### Custom Rules (Kodia)

| Rule | Description | Example |
|------|-------------|---------|
| `strong_password` | Requires uppercase, lowercase, digit, symbol | `validate:"strong_password"` |
| `phone` | Valid phone format (E.164 flexible) | `validate:"phone"` |
| `alpha_space` | Letters and spaces only | `validate:"alpha_space"` |
| `no_html` | No `<` or `>` characters | `validate:"no_html"` |
| `alphanumeric` | Letters and digits only | `validate:"alphanumeric"` |
| `slug` | URL slug format (lowercase, digits, hyphens) | `validate:"slug"` |
| `date_format:LAYOUT` | Date format check (Go time.Parse) | `validate:"date_format:2006-01-02"` |

#### Standard Rules (go-playground/validator)

| Rule | Description | Example |
|------|-------------|---------|
| `required` | Field must be non-empty | `validate:"required"` |
| `email` | Valid email address | `validate:"email"` |
| `min=N` | Minimum length/value | `validate:"min=8"` |
| `max=N` | Maximum length/value | `validate:"max=100"` |
| `len=N` | Exact length | `validate:"len=6"` |
| `url` | Valid URL | `validate:"url"` |
| `uuid4` | Valid UUID v4 | `validate:"uuid4"` |
| `omitempty` | Skip validation if empty | `validate:"omitempty,email"` |

### Async Validation (Context-Aware)

Async validation rules require context (e.g., database access). They run **after** sync validation passes and use the `ctx:` prefix in the `validate:` tag.

#### Example: Unique Email Check

```go
// In setup (e.g., DI container)
validate := validation.New()
validate.RegisterContextRule("unique_email", func(ctx context.Context, email string) error {
    exists, err := userRepo.ExistsByEmail(ctx, email)
    if err != nil {
        return err  // Database error
    }
    if exists {
        return errors.New("already registered")
    }
    return nil
})

// In DTO
type RegisterRequest struct {
    Email string `json:"email" validate:"required,email,ctx:unique_email"`
}

// In handler
if err := binder.BindCtx(c, &req); err != nil {
    return  // Async validation error sent as 422
}
```

#### Custom Async Rules

```go
// RegisterContextRule registers a rule that needs context
func (vl *Validator) RegisterContextRule(name string, fn ContextRule) {
    vl.ctxRules[name] = fn
}

// ContextRule is a validation function that requires context
type ContextRule func(ctx context.Context, fieldValue string) error
```

---

## Error Handling

### HTTP Status Codes

- **400 Bad Request** — Invalid JSON or type mismatch
- **422 Unprocessable Entity** — Validation rule failed

### Error Response Format

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

### Error Message Templates

Error messages are localized based on the `Accept-Language` header. Messages use the format:

```
{field} {message}
```

Example:
- English: `"email must be a valid email address"`
- Indonesian: `"email harus berupa alamat email yang valid"`

### Supported Locales

- `en` — English (default)
- `id` — Indonesian (Bahasa Indonesia)

### Extending Locales

```go
import "github.com/kodia-studio/kodia/pkg/validation"

// Register new locale or override existing
validation.RegisterLocale("id", validation.LocaleMessages{
    "required": "{field} wajib diisi",
    "email":    "{field} harus berupa email yang valid",
})
```

---

## Usage Guide

### Basic Pattern: Bind + Validate + Respond

```go
// 1. Define DTO with validation and sanitization rules
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=5,max=200" sanitize:"trim"`
    Content string `json:"content" validate:"required,min=50,max=5000" sanitize:"trim"`
    Tags    string `json:"tags" validate:"omitempty,max=100" sanitize:"trim,lowercase"`
}

// 2. In handler, bind + validate in one call
func (h *PostHandler) Create(c *gin.Context) {
    var req CreatePostRequest
    
    // Bind JSON → Sanitize → Validate
    if err := binder.Bind(c, &req); err != nil {
        return  // Error response already sent (400 or 422)
    }
    
    // req is guaranteed valid and sanitized here
    result, err := h.postService.Create(c.Request.Context(), req.Title, req.Content)
    
    if err != nil {
        response.InternalServerError(c, "Failed to create post")
        return
    }
    
    response.Created(c, "Post created", result)
}
```

### Pattern: Async Validation

```go
// DTO with async rule
type RegisterRequest struct {
    Email string `json:"email" validate:"required,email,ctx:unique_email" sanitize:"trim,lowercase"`
}

// In handler, use BindCtx for async validation
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    
    // Bind → Sanitize → Sync Validate → Async Validate
    if err := binder.BindCtx(c, &req); err != nil {
        return  // Error response already sent
    }
    
    // Process
    user, err := h.authService.Register(c.Request.Context(), req.Email)
    // ...
}
```

### Pattern: Custom Validation Message

```go
// In handler, handle specific validation errors
if err := binder.Bind(c, &req); err != nil {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        // Custom handling
        for _, fe := range ve {
            if fe.Tag() == "strong_password" {
                // Provide extra guidance
            }
        }
    }
    return
}
```

### Pattern: Middleware Validation

```go
// Use BindAndValidate in middleware for shared validation
type CreateCommentReq struct {
    Content string `json:"content" validate:"required,min=1,max=1000" sanitize:"trim"`
}

// Middleware validates before reaching handler
func validateCommentMW(vl *validation.Validator) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateCommentReq
        if !validation.BindAndValidate(c, vl, &req) {
            return  // Error response sent, abort
        }
        c.Set("comment_req", req)  // Store validated request
        c.Next()
    }
}
```

---

## Best Practices

### 1. Always Sanitize Before Validating

Sanitization removes whitespace and normalizes input. Always include sanitization for text fields:

```go
// ✓ Good
Name string `json:"name" sanitize:"trim" validate:"required,min=2"`

// ✗ Avoid
Name string `json:"name" validate:"required,min=2"`  // Whitespace passes validation
```

### 2. Use Strong Password Validation

All password fields must use `strong_password`:

```go
// ✓ Good
Password string `json:"password" validate:"required,min=8,max=72,strong_password"`

// ✗ Avoid
Password string `json:"password" validate:"required,min=8"`
```

### 3. Sanitize Emails to Lowercase

Email comparisons are case-insensitive, so normalize to lowercase:

```go
// ✓ Good
Email string `json:"email" sanitize:"trim,lowercase" validate:"required,email"`

// ✗ Avoid
Email string `json:"email" validate:"required,email"`
```

### 4. Use no_html for User-Generated Content

Prevent HTML/XSS injection:

```go
// ✓ Good
Title   string `json:"title" sanitize:"trim" validate:"required,no_html"`
Content string `json:"content" sanitize:"trim" validate:"required,no_html"`

// ✗ Avoid
Title   string `json:"title" validate:"required"`
```

### 5. Use alpha_space for Names

Prevent special characters in names:

```go
// ✓ Good
Name string `json:"name" sanitize:"trim" validate:"required,alpha_space"`

// ✗ Avoid
Name string `json:"name" validate:"required"`
```

### 6. Keep Validation in DTOs, Business Logic in Services

```go
// Handler layer: format validation only
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email" sanitize:"lowercase"`
    Password string `json:"password" validate:"required"`
}

// Service layer: business logic
func (s *AuthService) Login(ctx context.Context, email, password string) error {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil || user == nil {
        return errors.New("invalid credentials")  // Business rule
    }
    if !s.comparePassword(password, user.PasswordHash) {
        return errors.New("invalid credentials")  // Business rule
    }
    return nil
}
```

### 7. Register Async Rules in DI Setup

```go
// In your DI/container setup
func setupValidation(userRepo *repositories.UserRepository) *validation.Validator {
    vl := validation.New()
    
    // Register async rules
    vl.RegisterContextRule("unique_email", func(ctx context.Context, email string) error {
        exists, err := userRepo.ExistsByEmail(ctx, email)
        if err != nil {
            return err
        }
        if exists {
            return errors.New("already registered")
        }
        return nil
    })
    
    vl.RegisterContextRule("exists_user_id", func(ctx context.Context, userID string) error {
        user, err := userRepo.FindByID(ctx, userID)
        if err != nil {
            return err
        }
        if user == nil {
            return errors.New("user not found")
        }
        return nil
    })
    
    return vl
}
```

---

## Real-World Examples

### User Registration

```go
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100,alpha_space,no_html" sanitize:"trim"`
    Email    string `json:"email" validate:"required,email,ctx:unique_email" sanitize:"trim,lowercase"`
    Password string `json:"password" validate:"required,min=8,max=72,strong_password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := binder.BindCtx(c, &req); err != nil {
        return  // Validation + async check done
    }
    
    user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
    // ...
}
```

### Blog Post Creation

```go
type CreatePostRequest struct {
    Title       string   `json:"title" validate:"required,min=5,max=200" sanitize:"trim"`
    Slug        string   `json:"slug" validate:"required,slug" sanitize:"slug"`
    Description string   `json:"description" validate:"required,min=50,max=5000,no_html" sanitize:"trim"`
    Tags        []string `json:"tags" validate:"max=10,dive,min=1,max=50"`
}

func (h *PostHandler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := binder.Bind(c, &req); err != nil {
        return
    }
    
    post, err := h.postService.Create(c.Request.Context(), req)
    // ...
}
```

### Contact Form

```go
type ContactRequest struct {
    Name    string `json:"name" validate:"required,min=2,max=100,alpha_space" sanitize:"trim"`
    Email   string `json:"email" validate:"required,email" sanitize:"trim,lowercase"`
    Subject string `json:"subject" validate:"required,min=5,max=100,no_html" sanitize:"trim"`
    Message string `json:"message" validate:"required,min=10,max=5000,no_html" sanitize:"trim"`
}

func (h *ContactHandler) Send(c *gin.Context) {
    var req ContactRequest
    if err := binder.Bind(c, &req); err != nil {
        return
    }
    
    err := h.emailService.SendContact(c.Request.Context(), req.Email, req.Subject, req.Message)
    // ...
}
```

---

## Testing

### Unit Test Example

```go
func TestSanitizeTrim(t *testing.T) {
    type Req struct {
        Field string `sanitize:"trim"`
    }
    req := &Req{Field: "  hello  "}
    _ = validation.Sanitize(req)
    if req.Field != "hello" {
        t.Errorf("got %q, want %q", req.Field, "hello")
    }
}

func TestValidateEmail(t *testing.T) {
    type Req struct {
        Email string `validate:"required,email"`
    }
    
    vl := validation.New()
    req := &Req{Email: "invalid"}
    err := vl.Struct(req)
    if err == nil {
        t.Errorf("expected validation error for invalid email")
    }
}
```

### Integration Test Example

```go
func TestRegisterValidation(t *testing.T) {
    app := setupTestApp(t)
    
    tests := []struct {
        name    string
        body    interface{}
        wantCode int
    }{
        {
            name: "valid registration",
            body: map[string]string{
                "name":     "John Doe",
                "email":    "john@example.com",
                "password": "SecureP@ss123",
            },
            wantCode: http.StatusCreated,
        },
        {
            name: "invalid email",
            body: map[string]string{
                "name":     "John Doe",
                "email":    "invalid",
                "password": "SecureP@ss123",
            },
            wantCode: http.StatusUnprocessableEntity,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.body)
            req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            app.ServeHTTP(w, req)
            
            if w.Code != tt.wantCode {
                t.Errorf("got status %d, want %d", w.Code, tt.wantCode)
            }
        })
    }
}
```

---

## API Reference

### binder.Bind

```go
func (b *Binding) Bind(c *gin.Context, dest interface{}) error
```

Binds JSON → Sanitizes → Validates. Returns error if any step fails. Error response already sent to client.

```go
var req MyRequest
if err := binder.Bind(c, &req); err != nil {
    return  // 400 or 422 response sent
}
```

### binder.BindCtx

```go
func (b *Binding) BindCtx(c *gin.Context, dest interface{}) error
```

Same as `Bind` but includes async (context-aware) validation. Required for DTOs with `ctx:` rules.

```go
var req RegisterRequest
if err := binder.BindCtx(c, &req); err != nil {
    return  // 400, 422, or 422 (async) response sent
}
```

### validation.BindAndValidate

```go
func BindAndValidate(c *gin.Context, vl *Validator, req any) bool
```

Returns `true` if binding and validation succeed. Returns `false` and sends error response if they fail.

```go
var req MyRequest
if !validation.BindAndValidate(c, h.validator, &req) {
    return
}
```

### validation.BindAndValidateCtx

```go
func BindAndValidateCtx(c *gin.Context, vl *Validator, req any) bool
```

Same as `BindAndValidate` but includes async validation.

```go
var req RegisterRequest
if !validation.BindAndValidateCtx(c, h.validator, &req) {
    return
}
```

### validation.Sanitize

```go
func Sanitize(v any) error
```

Applies `sanitize:` tags to string fields. Used internally by `Bind` and middleware.

```go
var req MyRequest
_ = validation.Sanitize(&req)  // Error ignored (safe to skip)
```

### validation.FormatErrors

```go
func FormatErrors(err error) map[string][]string
```

Converts validation error to map of field → errors. Used for custom error handling.

```go
if err := vl.Struct(req); err != nil {
    errs := validation.FormatErrors(err)
    response.UnprocessableEntity(c, "Validation failed", errs)
}
```

### validation.FormatErrorsLocale

```go
func FormatErrorsLocale(err error, locale string) map[string][]string
```

Same as `FormatErrors` but uses specified locale for error messages.

```go
errs := validation.FormatErrorsLocale(err, "id")  // Indonesian
```

### Validator.RegisterContextRule

```go
func (vl *Validator) RegisterContextRule(name string, fn ContextRule)
```

Registers an async validation rule. Call during setup/DI initialization.

```go
vl.RegisterContextRule("unique_email", func(ctx context.Context, email string) error {
    exists, _ := repo.ExistsByEmail(ctx, email)
    if exists {
        return errors.New("already registered")
    }
    return nil
})
```

---

## Troubleshooting

### Validation Not Triggering

**Problem:** `validate:` tags are being ignored.

**Solution:** Make sure tags say `validate:` not `binding:`:

```go
// ✓ Correct
Email string `json:"email" validate:"required,email"`

// ✗ Wrong
Email string `json:"email" binding:"required,email"`
```

### Sanitization Not Applied

**Problem:** Input whitespace still present after sanitization.

**Solution:** Ensure `sanitize:` tag is present:

```go
// ✓ Correct
Name string `json:"name" sanitize:"trim" validate:"required"`

// ✗ Wrong
Name string `json:"name" validate:"required"`
```

### Custom Validation Rules Not Found

**Problem:** "validation not found" error.

**Solution:** Use `validation.New()` which registers custom rules:

```go
// ✓ Correct
vl := validation.New()

// ✗ Wrong
vl := validator.New()  // Standard go-playground validator
```

### Async Validation Returns 400 Instead of 422

**Problem:** Context-aware validation errors appear as 400.

**Solution:** Use `binder.BindCtx` or `validation.BindAndValidateCtx`:

```go
// ✓ Correct
if err := binder.BindCtx(c, &req); err != nil {
    return
}

// ✗ Wrong
if err := binder.Bind(c, &req); err != nil {
    return
}
```

---

## See Also

- [Validation Rules Reference](../pkg/validation/rules.go)
- [Sanitizer Implementation](../pkg/validation/sanitizer.go)
- [Validator Implementation](../pkg/validation/validator.go)
- [kodia-web Validation Integration](../../kodia-web/backend/docs/VALIDATION_LAYER.md)
