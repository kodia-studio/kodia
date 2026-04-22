# Real-time & Communication

Kodia provides a comprehensive, production-grade suite for real-time communication. It covers WebSocket, Server-Sent Events (SSE), unified notification delivery, and event broadcasting—all wired together out of the box.

---

## Architecture Overview

```
Domain Event System (EventDispatcher)
         │
         ▼
EventBroadcaster ────────────────────────────────────┐
         │                                           │
         ▼                                           ▼
   WebSocket Hub                             SSE Manager
         │                                           │
   ┌─────┴──────┐                         ┌──────────┘
   │  Channels  │                         │  SSE Clients
   │  Rooms     │                         │  (user/public)
   │  Users     │                         └──────────────
   └────────────┘

Notification Layer:
notifManager.Send(ctx, user, MyNotification{})
         │
   ┌─────┼──────┬──────┐
   ▼     ▼      ▼      ▼
  📧   📱     💬     🔔
Email  SMS   Slack  Push (FCM)
              + WebSocket (in-app)
```

---

## 1. WebSocket Hub

Kodia's WebSocket Hub supports rooms, named channels, user targeting, and presence tracking.

### Connecting a Client

```
ws://your-app.com/api/v1/ws          → general connection
ws://your-app.com/api/v1/ws/room/123 → join a specific room
```

### Hub API (from your services)

```go
// Get the hub from the app container
hub := app.MustGet("ws_hub").(*websocket.Hub)

// Broadcast to all clients
hub.Broadcast(&websocket.Message{
    Type:    websocket.MessageTypeBroadcast,
    Event:   "order.shipped",
    Payload: map[string]any{"order_id": "123"},
})

// Send to a specific user
hub.SendToUser(userID, &websocket.Message{
    Type:    websocket.MessageTypeNotification,
    Payload: websocket.NotificationPayload{Title: "Your order shipped!"},
})

// Send to a room
hub.SendToRoom("chat-room-1", &websocket.Message{
    Type:    websocket.MessageTypeChat,
    Payload: websocket.ChatPayload{Content: "Hello!"},
})

// Get online users
onlineUsers := hub.GetOnlineUsers()

// Get room presence
members := hub.GetRoomPresence("chat-room-1")

// Hub metrics
metrics := hub.Metrics()
// {"total_clients": 42, "total_users": 38, "total_rooms": 5, ...}
```

### Named Channels (Broadcasting)

```go
// Get or create a typed channel
ch := hub.GetOrCreateChannel("presence-team-42")

// Add middleware to the channel
ch.Use(func(conn *websocket.Connection, msg *websocket.Message) bool {
    // Only deliver to authenticated users
    return conn.UserID != ""
})

// Broadcast to the channel
hub.BroadcastToChannel("presence-team-42", &websocket.Message{
    Type:  websocket.MessageTypeBroadcast,
    Event: "team.update",
})
```

---

## 2. Server-Sent Events (SSE)

SSE is a lightweight alternative to WebSocket for **one-directional** data streams (server → client). It uses standard HTTP and works through firewalls and proxies.

### Endpoints

```
GET /api/v1/sse/:channel        → Subscribe to a public channel
GET /api/v1/sse/user            → Subscribe to private user stream (JWT required)
GET /api/v1/sse/status          → View active connection count
```

### Publishing from Your Backend

```go
// Get the SSE manager from the app container
sseManager := app.MustGet("sse_manager").(*sse.Manager)

// Publish to a channel
sseManager.Publish("orders", "order.shipped", map[string]any{
    "order_id": "abc-123",
    "status":   "shipped",
})

// Publish to a specific user
sseManager.PublishToUser(userID, "notification.new", payload)
```

### JavaScript Client Example

```javascript
const es = new EventSource('/api/v1/sse/orders');

es.addEventListener('order.shipped', (e) => {
    const data = JSON.parse(e.data);
    console.log('Order shipped:', data.order_id);
});

es.addEventListener('connected', () => {
    console.log('SSE connected!');
});

es.onerror = () => es.close();
```

---

## 3. Notification Channels

Kodia provides a unified interface to send notifications across multiple channels from a single `.Send()` call. Inspired by Laravel's notification system.

### Creating a Notification

```go
// Define your notification struct
type OrderShippedNotification struct {
    OrderID   string
    TrackingNumber string
}

func (n OrderShippedNotification) Via(notifiable ports.Notifiable) []string {
    // Return channels to use — any combination
    return []string{"email", "sms", "push", "websocket"}
}

func (n OrderShippedNotification) ToNotification(channel string, notifiable ports.Notifiable) *ports.NotificationMessage {
    switch channel {
    case "email":
        return &ports.NotificationMessage{
            Subject:  "Your order has shipped!",
            HtmlBody: fmt.Sprintf("<h1>Order %s shipped</h1>", n.OrderID),
        }
    case "sms":
        return &ports.NotificationMessage{
            SMSText: fmt.Sprintf("Your order %s has been shipped. Tracking: %s", n.OrderID, n.TrackingNumber),
        }
    case "push":
        return &ports.NotificationMessage{
            PushTitle: "Order Shipped!",
            PushBody:  "Your order is on its way.",
        }
    case "websocket":
        return &ports.NotificationMessage{
            WSEvent:   "order.shipped",
            WSPayload: map[string]any{"order_id": n.OrderID},
        }
    }
    return nil
}
```

### Sending a Notification

Your `User` model must implement `ports.Notifiable`:

```go
func (u *User) GetID() string          { return u.ID }
func (u *User) GetEmail() string       { return u.Email }
func (u *User) GetPhoneNumber() string { return u.PhoneNumber }
func (u *User) GetPushToken() string   { return u.FCMToken }
```

Then in your service:

```go
nm := app.MustGet("notification_manager").(*notification.Manager)

err := nm.Send(ctx, user, OrderShippedNotification{
    OrderID:        "abc-123",
    TrackingNumber: "1Z999AA10123456784",
})
```

### Available Channels

| Channel     | Driver         | Config Variable          |
|-------------|----------------|--------------------------|
| `email`     | Mailer (SMTP)  | `MAIL_*`                 |
| `websocket` | WS Hub         | Always active            |
| `sms`       | Twilio         | `NOTIFICATION_TWILIO_*`  |
| `slack`     | Webhook        | `NOTIFICATION_SLACK_*`   |
| `push`      | Firebase FCM   | `NOTIFICATION_FCM_*`     |

### Environment Variables

```env
# SMS (Twilio)
NOTIFICATION_TWILIO_ACCOUNT_SID=ACxxxxxxxx
NOTIFICATION_TWILIO_AUTH_TOKEN=your_token
NOTIFICATION_TWILIO_FROM_NUMBER=+1234567890

# Slack
NOTIFICATION_SLACK_WEBHOOK_URL=https://hooks.slack.com/services/xxx

# Push (Firebase)
NOTIFICATION_FCM_SERVER_KEY=AAAAxxxxxxx
```

---

## 4. Broadcast Events (Event → Real-time Bridge)

Integrate your domain event system with WebSocket and SSE automatically — like Laravel Broadcasting.

### Creating a Broadcastable Event

```go
// Implement ports.BroadcastEvent
type OrderStatusChanged struct {
    OrderID string
    Status  string
    UserID  string
}

func (e OrderStatusChanged) Name() string     { return "OrderStatusChanged" }
func (e OrderStatusChanged) Payload() interface{} { return e }

// BroadcastOn declares which channels receive this event
func (e OrderStatusChanged) BroadcastOn() []string {
    return []string{
        "public.orders",                          // All subscribers
        fmt.Sprintf("private-%s", e.UserID),     // Specific user's private channel
    }
}

// BroadcastAs overrides the event name sent to the client (optional)
func (e OrderStatusChanged) BroadcastAs() string { return "order.status.changed" }

// BroadcastWith defines the data sent to the client (optional)
func (e OrderStatusChanged) BroadcastWith() map[string]interface{} {
    return map[string]interface{}{
        "order_id": e.OrderID,
        "status":   e.Status,
    }
}
```

### Dispatching

```go
// Get the event broadcaster
broadcaster := app.MustGet("event_broadcaster").(*broadcasting.EventBroadcaster)

// Broadcast the event — it goes to WebSocket + SSE simultaneously
err := broadcaster.Broadcast(ctx, OrderStatusChanged{
    OrderID: "abc-123",
    Status:  "shipped",
    UserID:  "user-456",
})

// Or send ad-hoc events
broadcaster.BroadcastToUser(ctx, userID, "notification.new", data)
broadcaster.BroadcastToRoom(ctx, "team-42", "team.update", data)
```

### Channel Naming Convention

| Prefix        | Type     | Access          |
|---------------|----------|-----------------|
| `public.*`    | Public   | Anyone          |
| `private-{id}`| Private  | Authenticated   |
| `presence-{id}`| Presence | Auth + tracking |

---

## 5. Registering the Provider

Add `RealtimeProvider` to your app's provider list:

```go
app.Register(
    providers.NewWebSocketProvider(),   // Must be before RealtimeProvider
    providers.NewRealtimeProvider(),    // SSE + Broadcasting + Notifications
)
```
