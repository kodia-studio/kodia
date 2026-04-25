# Event System

Kodia Framework includes a complete event-driven architecture for handling domain events. Events enable loose coupling between different parts of your application and support both synchronous and asynchronous event handling via background jobs.

---

## Overview

The event system provides:

- **Domain events** — represent something that happened in your domain
- **Event listeners** — respond to events (sync or async)
- **Queue integration** — async listeners are automatically queued via Redis/Asynq
- **Built-in events** — pre-defined events for common domain operations
- **Developer-friendly API** — simple `Emit()` and `On()` methods

---

## Quick Start

### 1. Emit an Event

```go
import (
    "github.com/kodia-studio/kodia/internal/core/events"
    "github.com/kodia-studio/kodia/pkg/events"
)

// In a handler or service
func RegisterUser(c *gin.Context, userService ports.UserService, dispatcher *events.Dispatcher) {
    user, _ := userService.Register(ctx, input)

    // Emit event
    dispatcher.Emit(ctx, events.UserRegisteredEvent{
        UserID:    user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Timestamp: time.Now(),
    })
}
```

### 2. Listen to Events (Synchronous)

```go
import (
    "github.com/kodia-studio/kodia/internal/core/events"
    "github.com/kodia-studio/kodia/internal/core/ports"
)

// Define a listener
type WelcomeEmailListener struct {
    mailer ports.MailProvider
}

func (l *WelcomeEmailListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(events.UserRegisteredEvent)
    
    // Send welcome email
    return l.mailer.Send(ctx, &mail.Mail{
        To:      e.Email,
        Subject: "Welcome!",
        Body:    fmt.Sprintf("Hello %s, welcome to Kodia!", e.Name),
    })
}

// Register listener in your bootstrap
dispatcher.On("user.registered", &WelcomeEmailListener{mailer})
```

### 3. Async Event Handling (Queue-based)

```go
// Define async listener
type SendAnalyticsListener struct {
    analytics analytics.Service
}

func (l *SendAnalyticsListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(events.UserRegisteredEvent)
    return l.analytics.TrackUserSignup(ctx, e.UserID, e.Email)
}

// Mark as async by implementing ShouldQueue
func (l *SendAnalyticsListener) ShouldQueue() bool {
    return true  // This listener will run in background
}

// Register listener
dispatcher.On("user.registered", &SendAnalyticsListener{analytics})
```

When `ShouldQueue()` returns `true`, the event is serialized to Redis and processed asynchronously by a worker process. This keeps your request-response cycle fast while still handling the event.

---

## Built-in Domain Events

Kodia provides common domain events ready to use. All events are in `internal/core/events/built_in.go`:

### User Events

```go
// Dispatched on new user registration
events.UserRegisteredEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Name:      "John",
    Timestamp: time.Now(),
}

// Dispatched on successful login
events.UserLoggedInEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    IPAddress: "192.168.1.1",
    UserAgent: "Mozilla/5.0...",
    Timestamp: time.Now(),
}
```

### Password Events

```go
// Dispatched when password is changed
events.PasswordChangedEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Timestamp: time.Now(),
}

// Dispatched when password reset is requested/completed
events.PasswordResetEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Status:    "completed",  // "requested" or "completed"
    Timestamp: time.Now(),
}
```

### 2FA Events

```go
// Dispatched when 2FA is enabled
events.TwoFactorEnabledEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Timestamp: time.Now(),
}

// Dispatched when 2FA is disabled
events.TwoFactorDisabledEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Timestamp: time.Now(),
}
```

### Email Events

```go
// Dispatched when email verification is requested
events.EmailVerificationRequestedEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Timestamp: time.Now(),
}

// Dispatched when email is verified
events.EmailVerifiedEvent{
    UserID:    "user123",
    Email:     "user@example.com",
    Timestamp: time.Now(),
}
```

### Role Events

```go
// Dispatched when a role is assigned to a user
events.RoleAssignedEvent{
    UserID:    "user123",
    RoleName:  "editor",
    Timestamp: time.Now(),
}

// Dispatched when a role is revoked
events.RoleRevokedEvent{
    UserID:    "user123",
    RoleName:  "editor",
    Timestamp: time.Now(),
}
```

---

## Creating Custom Events

### 1. Define Your Event

Events must implement the `ports.Event` interface:

```go
package events

type OrderCreatedEvent struct {
    OrderID   string
    UserID    string
    Amount    float64
    Items     int
    Timestamp time.Time
}

func (e OrderCreatedEvent) Name() string {
    return "order.created"
}

func (e OrderCreatedEvent) Payload() interface{} {
    return e
}
```

Or use the convenience `BaseEvent`:

```go
import "github.com/kodia-studio/kodia/pkg/events"

// Simple event
event := events.NewEvent("custom.event", map[string]interface{}{
    "order_id": "123",
    "amount": 99.99,
})

dispatcher.Emit(ctx, event)
```

### 2. Define Listeners

```go
// Sync listener
type SendOrderConfirmationListener struct {
    mailer ports.MailProvider
}

func (l *SendOrderConfirmationListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(OrderCreatedEvent)
    // Send confirmation email
    return l.mailer.Send(...)
}

// Async listener (runs in background)
type UpdateInventoryListener struct {
    inventory InventoryService
}

func (l *UpdateInventoryListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(OrderCreatedEvent)
    return l.inventory.ReserveItems(ctx, e.OrderID, e.Items)
}

func (l *UpdateInventoryListener) ShouldQueue() bool {
    return true
}
```

### 3. Register Listeners

```go
dispatcher := app.Resolve[*events.Dispatcher]("events")

dispatcher.On("order.created",
    &SendOrderConfirmationListener{mailer},
    &UpdateInventoryListener{inventory},
)
```

### 4. Emit the Event

```go
func CreateOrder(c *gin.Context, orderService OrderService, dispatcher *events.Dispatcher) {
    order, _ := orderService.Create(ctx, input)

    dispatcher.Emit(ctx, OrderCreatedEvent{
        OrderID:   order.ID,
        UserID:    order.UserID,
        Amount:    order.Total,
        Items:     len(order.Items),
        Timestamp: time.Now(),
    })
}
```

---

## CLI Generators

The Kodia CLI includes generators for events and listeners:

### Generate Event

```bash
kodia make:event OrderCreated
# Creates: internal/core/events/order_created_event.go
```

### Generate Listener

```bash
kodia make:listener SendOrderConfirmation
# Creates: internal/core/listeners/send_order_confirmation_listener.go
# Auto-registers in: internal/core/events/registry.go
```

---

## Best Practices

✅ **Do:**
- Use events to decouple domains (instead of direct service calls)
- Mark time-consuming listeners as async (`ShouldQueue() = true`)
- Include relevant context in event payloads (IDs, names, timestamps)
- Keep listeners focused on a single concern
- Name events in past tense (`UserRegistered`, `OrderCreated`)
- Use descriptive event names with dot notation (`order.created`, `payment.processed`)
- Implement proper error handling in listeners

❌ **Don't:**
- Make sync listeners too slow (move to async if > 100ms)
- Emit events from listeners (can cause infinite loops)
- Store large objects in event payloads (IDs + necessary data only)
- Forget to handle listener failures gracefully
- Create overly granular events for every tiny change

---

## Architecture

### Internal Dispatcher (`internal/infrastructure/events/dispatcher.go`)

The core engine that routes events to listeners:

```go
type InternalDispatcher struct {
    listeners map[string][]ports.Listener
    queue     ports.QueueProvider  // Redis/Asynq
    log       *zap.Logger
}

// Dispatch sends event to listeners
// Sync listeners: execute immediately
// Async listeners: queue to Redis
func (d *InternalDispatcher) Dispatch(ctx, event) error { ... }
```

### Public Dispatcher (`pkg/events/dispatcher.go`)

Developer-facing API wrapping the internal dispatcher:

```go
type Dispatcher struct {
    internal ports.EventDispatcher
}

func (d *Dispatcher) Emit(ctx, event) error
func (d *Dispatcher) On(eventName string, listeners ...)
```

### Registry (`internal/core/events/registry.go`)

Bootstrap function that registers all listeners:

```go
func RegisterEvents(dispatcher ports.EventDispatcher) {
    dispatcher.On("user.registered", &SendWelcomeEmailListener{...})
    dispatcher.On("order.created", &SendConfirmationListener{...})
    // ... more listeners
}
```

---

## Sync vs Async Comparison

| Aspect | Sync | Async |
|---|---|---|
| **Execution** | Immediate, blocking | Queued, non-blocking |
| **Best for** | Fast, simple operations | Slow, fire-and-forget tasks |
| **Example** | Validation, audit logging | Email, analytics, heavy processing |
| **Error handling** | Return errors to caller | Retry via queue |
| **Performance** | Slower responses if handlers are slow | Faster responses |
| **Reliability** | Fails if listener fails | Retries via Redis |

**Rule of thumb:** If a listener takes > 100ms, make it async.

---

## Error Handling

Listener errors are logged but don't stop other listeners:

```go
// Internal dispatcher logs and continues
for _, listener := range listeners {
    if err := listener.Handle(ctx, event); err != nil {
        dispatcher.log.Error("Listener error",
            zap.String("event", event.Name()),
            zap.Error(err),
        )
        // Continue to next listener
    }
}
```

For async listeners, Asynq automatically retries failed tasks with exponential backoff.

---

## Real-World Example: Order Flow

```go
// Handler: Create order
func (h *OrderHandler) Create(c *gin.Context) {
    orderService := h.app.Resolve[OrderService]("order_service")
    dispatcher := h.app.Resolve[*events.Dispatcher]("events")

    order, _ := orderService.Create(ctx, input)

    // Emit event (automatically triggers all registered listeners)
    dispatcher.Emit(ctx, events.OrderCreatedEvent{
        OrderID:   order.ID,
        UserID:    order.UserID,
        Amount:    order.Total,
        Items:     len(order.Items),
        Timestamp: time.Now(),
    })

    c.JSON(200, order)
}

// Listener 1: Send confirmation email (async)
type SendOrderConfirmationListener struct {
    mailer ports.MailProvider
}

func (l *SendOrderConfirmationListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(events.OrderCreatedEvent)
    return l.mailer.SendWithTemplate(ctx, &mail.Mail{
        To:       e.UserID,
        Template: "order_confirmation",
        Data: map[string]interface{}{
            "order_id": e.OrderID,
            "amount":   e.Amount,
        },
    })
}

func (l *SendOrderConfirmationListener) ShouldQueue() bool {
    return true  // Async
}

// Listener 2: Update inventory (async)
type UpdateInventoryListener struct {
    inventory InventoryService
}

func (l *UpdateInventoryListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(events.OrderCreatedEvent)
    return l.inventory.ReserveItems(ctx, e.OrderID, e.Items)
}

func (l *UpdateInventoryListener) ShouldQueue() bool {
    return true  // Async
}

// Listener 3: Audit log (sync)
type AuditLogListener struct {
    auditRepo AuditRepository
}

func (l *AuditLogListener) Handle(ctx context.Context, event ports.Event) error {
    e := event.(events.OrderCreatedEvent)
    return l.auditRepo.Log(ctx, "order_created", e.OrderID)
}

func (l *AuditLogListener) ShouldQueue() bool {
    return false  // Sync (very fast)
}
```

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
