# Sentry Integration — Error Monitoring & Performance Tracing

Complete guide to Kodia Framework's Sentry integration for production error monitoring, crash reporting, and performance tracing.

---

## Overview

Sentry is integrated to provide:
- **Error Monitoring** — automatic panic capture and error reporting
- **Performance Tracing** — track response times per endpoint
- **Release Tracking** — group issues by release version
- **User Context** — attach user information to errors
- **Breadcrumbs** — track events leading up to errors

---

## Setup

### 1. Environment Variables

Set your Sentry DSN in `.env`:

```bash
SENTRY_DSN=https://your-key@sentry.io/project-id
```

Optionally configure:

```bash
OBSERVABILITY_SAMPLING_RATE=0.1          # Trace sampling: 0.0-1.0 (default: 0.1)
OBSERVABILITY_TRACING_ENABLED=true       # Enable OpenTelemetry tracing
OBSERVABILITY_SERVICE_NAME=kodia-backend # Service name in Sentry
```

### 2. No Code Changes Required

Sentry is automatically initialized when the app starts if `SENTRY_DSN` is set. Panic recovery and HTTP performance tracing are built-in.

---

## How It Works

### Panic Recovery

Every panic is automatically caught by the Recovery middleware and reported to Sentry:

```go
// Any panic in a handler is automatically caught
func (h *Handler) DeleteUser(c *gin.Context) {
    panic("oops!")  // → Captured by Recovery → Reported to Sentry
}
```

### HTTP Performance Tracing

Every HTTP request creates a Sentry transaction tracking performance:

```
GET /api/users/123
  ├─ http.method: GET
  ├─ http.status: 200
  ├─ duration: 145ms
  └─ [Any child spans from your code]
```

Transactions are named by route: `GET /api/users/{id}`, `POST /api/posts`, etc.

---

## Developer API

Use `pkg/monitoring` helpers in your handlers and services:

### Capture Errors

```go
import "github.com/kodia-studio/kodia/pkg/monitoring"

func (h *Handler) CreatePost(c *gin.Context) {
    post, err := h.service.Create(ctx, input)
    if err != nil {
        // Report error with context tags
        monitoring.CaptureError(ctx, err, map[string]string{
            "user_id": c.GetString("user_id"),
            "action":  "create_post",
        })
        response.InternalServerError(c, "Failed to create post")
        return
    }
}
```

### Send Messages

```go
// Log informational messages
monitoring.CaptureMessage(ctx, "User upgraded to premium", sentry.LevelInfo)

// Log warnings
monitoring.CaptureMessage(ctx, "High error rate detected", sentry.LevelWarning)
```

### Attach User Context

```go
// In your auth middleware, after validating the user
type SentryUser struct {
    ID        string
    Email     string
    Username  string
    IPAddress string
}

monitoring.SetUser(ctx, monitoring.SentryUser{
    ID:       user.ID,
    Email:    user.Email,
    Username: user.Username,
    IPAddress: c.ClientIP(),
})
```

### Add Breadcrumbs

```go
// Track events that led to an error
monitoring.AddBreadcrumb(ctx, "database", "User query executed", sentry.LevelInfo, map[string]interface{}{
    "query": "SELECT * FROM users WHERE id = ?",
    "duration_ms": 42,
})

// Later if an error occurs:
// The breadcrumbs are automatically attached to the error in Sentry
```

### Create Spans

```go
// Wrap slow operations in a Sentry span for performance profiling
err := monitoring.WithSpan(ctx, "db.query", "Fetch user from database", func(ctx context.Context) error {
    user, err := userService.GetByID(ctx, userID)
    return err
})
```

### Check Sentry Status

```go
if monitoring.IsInitialized() {
    // Sentry is configured and ready
}
```

---

## Configuration Reference

| Environment Variable | Default | Description |
|---|---|---|
| `SENTRY_DSN` | — | Sentry DSN (required to enable Sentry) |
| `OBSERVABILITY_SAMPLING_RATE` | 0.1 | Trace sampling rate (0.0-1.0) |
| `OBSERVABILITY_TRACING_ENABLED` | true | Enable OpenTelemetry tracing |
| `OBSERVABILITY_SERVICE_NAME` | kodia | Service name reported to Sentry |

---

## What Gets Reported

### Automatically Captured

- ✅ Panics in HTTP handlers (via Recovery middleware)
- ✅ All HTTP request performance (via Tracing middleware)
- ✅ HTTP status codes and response times
- ✅ Request URL, method, and route pattern
- ✅ Stack traces with source code context
- ✅ User information (when `SetUser()` is called)
- ✅ Breadcrumbs (when `AddBreadcrumb()` is called)

### Manually Captured

- `CaptureError()` — report errors with context
- `CaptureMessage()` — send info/warning messages
- `SetUser()` — attach user context
- `AddBreadcrumb()` — track events
- `WithSpan()` — measure operation performance

---

## Real-World Example

```go
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
    ctx := c.Request.Context()
    userID := c.GetString("user_id")

    // Attach user context to all Sentry reports
    monitoring.SetUser(ctx, monitoring.SentryUser{
        ID:        userID,
        Email:     c.GetString("user_email"),
        Username:  c.GetString("user_name"),
        IPAddress: c.ClientIP(),
    })

    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        monitoring.CaptureError(ctx, err, map[string]string{
            "stage": "validation",
        })
        response.BadRequest(c, err.Error())
        return
    }

    // Add breadcrumb for auditing
    monitoring.AddBreadcrumb(ctx, "user", "Creating new user", sentry.LevelInfo, map[string]interface{}{
        "email": input.Email,
    })

    // Wrap database operation in a span
    var user *domain.User
    err := monitoring.WithSpan(ctx, "db.create", "Create user in database", func(ctx context.Context) error {
        var err error
        user, err = h.userService.Create(ctx, input)
        return err
    })

    if err != nil {
        monitoring.CaptureError(ctx, err, map[string]string{
            "stage": "create_user",
            "email": input.Email,
        })
        response.InternalServerError(c, "Failed to create user")
        return
    }

    response.Created(c, "User created", user)
}
```

In Sentry, this will appear as:

```
Error: <error message>
  Level: error
  User: John Doe <john@example.com>
  Breadcrumbs:
    - [user] Creating new user (email: john@example.com)
  Spans:
    - db.create (Fetch user from database) — 145ms
    - [Parent] POST /api/users — 156ms
```

---

## Performance Tracing

### Transaction View

In Sentry, view endpoint performance:

```
Transactions
  GET /api/users          avg: 145ms     (1000 samples)
  POST /api/users         avg: 234ms     (450 samples)
  GET /api/users/{id}     avg: 42ms      (5000 samples)
  POST /api/posts         avg: 567ms     (200 samples)
```

### Custom Spans

Create nested spans for detailed profiling:

```go
err := monitoring.WithSpan(ctx, "db.query.users", "SELECT * FROM users", func(ctx context.Context) error {
    // Query runs here
    return userRepo.FindAll(ctx)
})

// In Sentry:
// POST /api/admin/users (534ms)
//   └─ db.query.users (287ms)
//   └─ auth.validate (42ms)
//   └─ response.marshal (12ms)
```

---

## Best Practices

✅ **Do:**
- Set user context in your auth middleware for all errors to be tagged with user info
- Use `AddBreadcrumb()` for important business logic checkpoints
- Wrap expensive operations in `WithSpan()` for performance tracking
- Use appropriate log levels (Debug, Info, Warning, Error, Fatal)
- Check `IsInitialized()` before calling Sentry helpers if Sentry might not be configured
- Set a meaningful `OBSERVABILITY_SERVICE_NAME` for multi-service deployments

❌ **Don't:**
- Report sensitive data (passwords, tokens, API keys) in breadcrumbs or error tags
- Create too many spans per request (adds overhead)
- Set `OBSERVABILITY_SAMPLING_RATE` to 1.0 in production (costs money, creates noise)
- Forget to call `SetUser()` in your auth middleware
- Rely only on Sentry — use it alongside structured logging (zap)

---

## Troubleshooting

### No events in Sentry

1. Verify `SENTRY_DSN` is set and correct
2. Check that `OBSERVABILITY_SAMPLING_RATE > 0`
3. Verify Sentry is initialized by checking logs: `"Sentry initialized successfully"`
4. Test by triggering a panic in a handler

### Missing user context

Call `SetUser()` in your auth middleware:

```go
// In your Auth middleware
monitoring.SetUser(c.Request.Context(), monitoring.SentryUser{
    ID:       user.ID,
    Email:    user.Email,
    Username: user.Name,
})
```

### Too many transactions

Lower `OBSERVABILITY_SAMPLING_RATE` to 0.1 or 0.01 in production.

### Performance overhead

Each Sentry transaction has minimal overhead (< 1ms per request). If performance is a concern:
1. Lower the sampling rate
2. Disable OTEL tracing if not needed (`OBSERVABILITY_TRACING_ENABLED=false`)

---

## Resources

- [Sentry Documentation](https://docs.sentry.io/)
- [Sentry Go SDK](https://github.com/getsentry/sentry-go)
- [Performance Monitoring Guide](https://docs.sentry.io/product/performance/)

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
