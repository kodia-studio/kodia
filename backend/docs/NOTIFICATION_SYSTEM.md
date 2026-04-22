# Notification System — Kodia Framework

## Overview

The Notification System provides a complete solution for sending, storing, and managing notifications in Kodia Framework. It features:

- **Real-time push** via WebSocket (instant user notification)
- **Persistent storage** in PostgreSQL (notification history)
- **Async email** notifications via Event/Listener system with Asynq
- **Pagination & filtering** for notification lists
- **Unread count** tracking

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│ Service Layer (Core Logic)                              │
│  • Create notification                                  │
│  • Persist to database                                  │
│  • Push real-time via Broadcaster                       │
│  • Dispatch email event for async processing            │
└──────────────────┬──────────────────────────────────────┘
                   │
      ┌────────────┼────────────┐
      │            │            │
      ▼            ▼            ▼
   Repository   Broadcaster   EventDispatcher
   (GORM)       (WebSocket)    (Event System)
      │            │            │
      ▼            ▼            ▼
   Database   Connected      Async Workers
              Clients        (Email, etc.)
```

## Components

### Domain Entity (`domain/entities.go`)

```go
type NotificationType string

const (
    NotificationTypeInfo    = "info"
    NotificationTypeSuccess = "success"
    NotificationTypeWarning = "warning"
    NotificationTypeError   = "error"
)

type Notification struct {
    ID        string
    UserID    string
    Type      NotificationType
    Title     string
    Message   string
    Data      map[string]interface{}
    IsRead    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Repository Interface (`ports/repositories.go`)

- `Create(ctx, notification)` — Save to database
- `FindByID(ctx, id)` — Fetch single notification
- `FindByUserID(ctx, userID, params)` — Paginated list for user
- `MarkAsRead(ctx, id, userID)` — Mark notification as read
- `MarkAllAsRead(ctx, userID)` — Mark all user notifications as read
- `Delete(ctx, id, userID)` — Remove notification
- `CountUnread(ctx, userID)` — Get unread count

### Service Interface (`ports/services.go`)

```go
type NotificationService interface {
    Send(ctx, input) (*Notification, error)
    GetAll(ctx, userID, params) ([]*Notification, int64, error)
    MarkAsRead(ctx, id, userID) error
    MarkAllAsRead(ctx, userID) error
    Delete(ctx, id, userID) error
    CountUnread(ctx, userID) (int64, error)
}

type SendNotificationInput struct {
    UserID    string
    Type      NotificationType
    Title     string
    Message   string
    Data      map[string]interface{}
    SendEmail bool  // Trigger async email notification
}
```

### HTTP Endpoints

#### List Notifications

```
GET /api/notifications?page=1&per_page=15
```

**Response:**
```json
{
  "success": true,
  "message": "Notifications retrieved",
  "data": {
    "notifications": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "info",
        "title": "Welcome!",
        "message": "Welcome to our platform",
        "data": null,
        "is_read": false,
        "created_at": "2024-04-19T10:00:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "per_page": 15
  }
}
```

#### Get Unread Count

```
GET /api/notifications/unread-count
```

**Response:**
```json
{
  "success": true,
  "message": "Unread count retrieved",
  "data": {
    "count": 5
  }
}
```

#### Mark as Read

```
PUT /api/notifications/{id}/read
```

#### Mark All as Read

```
PUT /api/notifications/read-all
```

#### Delete Notification

```
DELETE /api/notifications/{id}
```

## Integration Guide

### 1. Sending a Notification from a Service

```go
// In your service (e.g., OrderService)
notifService := app.MustGet("notification_service").(ports.NotificationService)

notif, err := notifService.Send(ctx, ports.SendNotificationInput{
    UserID:    orderUserID,
    Type:      domain.NotificationTypeSuccess,
    Title:     "Order Placed",
    Message:   "Your order #12345 has been placed successfully",
    Data: map[string]interface{}{
        "order_id": "12345",
        "total":    199.99,
    },
    SendEmail: true,  // Also send email notification
})
if err != nil {
    return err
}
```

### 2. Listening to Notifications via WebSocket

```javascript
// In your frontend (Svelte/JavaScript)
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    
    if (message.type === 'notification') {
        // Display notification to user
        console.log(message.payload);
        // payload: { title, message, data }
    }
};
```

### 3. Custom Email Listener

The system automatically sends emails for notifications with `SendEmail: true` via the `SendNotificationEmail` listener. To customize the email template, modify `resources/mail/notification.html`:

```html
<!-- resources/mail/notification.html -->
<h1>{{ .title }}</h1>
<p>{{ .message }}</p>
<div>{{ .data }}</div>
```

## Database Schema

```sql
CREATE TABLE notifications (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
);
```

## Best Practices

### ✓ Do

- Use appropriate `NotificationType` for each notification (info, success, warning, error)
- Include `Data` for dynamic content (order ID, product name, etc.)
- Set `SendEmail: true` only when an email is essential
- Paginate large notification lists
- Mark notifications as read when user views them

### ✗ Don't

- Send notifications for every minor action (create unnecessary noise)
- Store sensitive data in `Data` (use encryption if needed)
- Dispatch events synchronously in critical paths
- Forget to check `Authorize(userID)` before deleting

## Performance Considerations

- **Real-time push** is instant (WebSocket)
- **Database persistence** is optimized with indexes on `user_id` and `created_at`
- **Email sending** is async (Asynq worker) — doesn't block the HTTP response
- **Pagination** defaults to 15 per page with max 100 per page

## Testing

### Unit Test Example

```go
func TestNotificationService_Send(t *testing.T) {
    repo := &mockNotificationRepository{}
    broadcaster := &mockBroadcaster{}
    dispatcher := &mockEventDispatcher{}
    
    service := services.NewNotificationService(repo, broadcaster, dispatcher, log)
    
    notif, err := service.Send(ctx, ports.SendNotificationInput{
        UserID:  "user-123",
        Type:    domain.NotificationTypeInfo,
        Title:   "Test",
        Message: "Test message",
    })
    
    require.NoError(t, err)
    require.Equal(t, "Test", notif.Title)
    require.True(t, repo.CreateCalled)
    require.True(t, broadcaster.NotifyUserCalled)
}
```

## Troubleshooting

### Notifications not being saved

- Ensure `notification_provider.go` is registered in `main.go`
- Check database migrations have run (notifications table exists)
- Verify user ID is valid

### WebSocket push not received

- Confirm user is connected to WebSocket
- Check `Broadcaster` is properly wired in container
- Verify user ID matches in both service and WebSocket client

### Emails not being sent

- Ensure `SendEmail: true` is set
- Check Asynq worker is running
- Verify email template exists in `resources/mail/`

## See Also

- [WebSocket Guide](WEBSOCKET_GUIDE.md) — Real-time communication
- [Event System](TODO) — Async event handling
- [Email Configuration](TODO) — Mail setup
