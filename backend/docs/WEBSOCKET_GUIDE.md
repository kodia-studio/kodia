# 🔄 WebSocket & Real-time Features Guide

Complete guide to using WebSocket in Kodia Framework for real-time notifications, chat, and live updates.

**Table of Contents:**
- [What is WebSocket?](#what-is-websocket)
- [Architecture](#architecture)
- [Client Connection](#client-connection)
- [Message Types](#message-types)
- [Broadcasting Patterns](#broadcasting-patterns)
- [Room-Based Chat](#room-based-chat)
- [Notifications](#notifications)
- [Security](#security)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## What is WebSocket?

WebSocket provides **full-duplex bidirectional communication** over a single TCP connection:
- Persistent connection (unlike HTTP request-response)
- Low-latency messaging (no handshake per message)
- Perfect for: chat, notifications, live dashboards, collaborative tools

**vs HTTP Polling:**
- WebSocket: Server pushes data → low latency, low CPU
- Polling: Client asks repeatedly → high latency, wastes bandwidth

Kodia uses `gorilla/websocket` with a Hub-based architecture.

---

## Architecture

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │ GET /api/ws?token=<jwt>
       ▼
┌─────────────────────────┐
│  HTTP Upgrade Handler   │ (router.go)
│  - Validates JWT        │
│  - Upgrades to WS       │
└──────┬──────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│  Hub (Central Connection Manager)   │
│  - Maintains client registry        │
│  - Routes messages globally/room    │
│  - Manages user connections         │
└──────┬──────────────────────────────┘
       │
       ├─► Connection 1 (User A)
       │    ├─ readPump (receives)
       │    └─ writePump (sends)
       │
       ├─► Connection 2 (User B, Room "chat-1")
       │    ├─ readPump
       │    └─ writePump
       │
       └─► Connection N (...)
```

**Flow:**
1. Client connects: `GET /api/ws?token=<jwt>`
2. Handler validates JWT, upgrades to WebSocket
3. Connection registered in Hub
4. Messages flow: Client → readPump → Hub → writePump → Client

---

## Client Connection

### JavaScript / SvelteKit

```javascript
// Using the provided Svelte store (recommended)
import { wsStore } from '$lib/stores/websocket'

// Connect
wsStore.connect('your-jwt-token')

// Listen for messages
wsStore.onMessage('notification', (msg) => {
    console.log('Got notification:', msg)
})

// Send message
wsStore.send({
    type: 'chat',
    payload: { content: 'Hello!' }
})

// Disconnect
wsStore.disconnect()
```

### Browser Console (Testing)

```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws?token=eyJhbGc...')

ws.onopen = () => {
    console.log('Connected!')
    ws.send(JSON.stringify({
        type: 'chat',
        payload: { content: 'Hello WebSocket!' }
    }))
}

ws.onmessage = (event) => {
    console.log('Message:', JSON.parse(event.data))
}

ws.onerror = (err) => console.error('Error:', err)
ws.onclose = () => console.log('Disconnected')
```

---

## Message Types

Kodia defines standardized message types:

### Message Structure

```json
{
  "type": "notification",
  "payload": { /* varies by type */ },
  "room_id": "chat-room-1",
  "user_id": "user-123",
  "timestamp": 1713607200
}
```

### Notification

```json
{
  "type": "notification",
  "payload": {
    "title": "Order Shipped",
    "message": "Your order #123 is on the way",
    "data": { "order_id": "123", "tracking": "..." }
  }
}
```

### Chat Message

```json
{
  "type": "chat",
  "payload": {
    "sender_id": "user-456",
    "sender_name": "Alice",
    "content": "Hey everyone!",
    "room_id": "chat-1"
  }
}
```

### Ping/Pong (Heartbeat)

```json
{
  "type": "ping",
  "timestamp": 1713607200
}
```

Server responds with `type: "pong"`.

### Error

```json
{
  "type": "error",
  "payload": {
    "code": "UNAUTHORIZED",
    "message": "Token expired"
  }
}
```

### Custom Types

Add your own message types by extending the `MessageType` enum in `adapter/websocket/message.go`:

```go
const (
    MessageTypeCustomEvent MessageType = "custom_event"
)
```

---

## Broadcasting Patterns

### Global Broadcast (All Clients)

```go
// In a service
broadcaster.Broadcast("status_update", gin.H{
    "status": "system_maintenance",
    "duration": "2 hours",
})

// All connected clients receive this message
```

Usage: System announcements, global status changes.

### Send to Specific User

```go
broadcaster.NotifyUser("user-123", "order_update", gin.H{
    "order_id": "456",
    "status": "shipped",
})

// Only user-123's connections receive this
```

Usage: Personal notifications, account updates.

### Send to Multiple Users

```go
broadcaster.NotifyMultipleUsers(
    []string{"user-1", "user-2", "user-3"},
    "team_announcement",
    gin.H{
        "message": "Sprint planning at 3pm",
    },
)
```

### Broadcast to Room

```go
broadcaster.BroadcastToRoom("chat-room-1", "chat", gin.H{
    "sender_id": "user-456",
    "content": "Hello room!",
})

// All clients in "chat-room-1" receive this
```

---

## Room-Based Chat

### Connect to a Room

```javascript
const ws = new WebSocket(
  'ws://localhost:8080/api/ws/room/chat-room-1?token=<jwt>'
)

// Send message to room
ws.send(JSON.stringify({
  type: 'chat',
  payload: {
    content: 'Hello everyone in the room!'
  }
}))
```

### Server-Side (Broadcasting to Room)

```go
// In a chat service
func (s *ChatService) SendMessage(ctx context.Context, roomID string, userID string, content string) error {
    // Save to database
    msg := &domain.ChatMessage{
        RoomID: roomID,
        UserID: userID,
        Content: content,
    }
    if err := s.repo.Create(ctx, msg); err != nil {
        return err
    }

    // Broadcast to room in real-time
    s.broadcaster.BroadcastToRoom(roomID, "chat", gin.H{
        "sender_id": userID,
        "sender_name": "Alice",
        "content": content,
        "timestamp": time.Now().Unix(),
    })

    return nil
}
```

---

## Notifications

### Example: Order Status Updates

```go
// In OrderService
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
    order := &domain.Order{ ID: orderID, Status: status }
    
    if err := s.repo.Update(ctx, order); err != nil {
        return err
    }

    // Notify the user
    s.broadcaster.NotifyUser(order.UserID, "order_status_update", gin.H{
        "order_id": orderID,
        "status": status,
        "updated_at": time.Now(),
    })

    return nil
}
```

### Example: Multiple User Notification

```go
// Notify team members about task assignment
func (s *TaskService) AssignTask(ctx context.Context, taskID string, assigneeIDs []string) error {
    // ... update database ...

    s.broadcaster.NotifyMultipleUsers(
        assigneeIDs,
        "task_assigned",
        gin.H{
            "task_id": taskID,
            "title": "Review PR #123",
        },
    )

    return nil
}
```

---

## Security

### JWT Validation

Every WebSocket connection requires a valid JWT token:

```javascript
// ✅ Correct - token in query param or header
ws = new WebSocket('ws://localhost:8080/api/ws?token=<jwt>')

// ✅ Also works - Authorization header
ws = new WebSocket('ws://localhost:8080/api/ws')
// Set header: Authorization: Bearer <jwt>
```

The handler validates the JWT before upgrading:

```go
claims, err := h.jwtManager.ValidateAccessToken(token)
if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
    return
}
```

**Tokens expire:** If a client's token expires, they're disconnected on next message.

### CORS & Origin Checking

The WebSocket upgrader allows all origins by default (for development):

```go
upgrader.CheckOrigin = func(r *http.Request) bool {
    return true // In production, validate origin
}
```

**For production**, implement proper origin checking:

```go
upgrader.CheckOrigin = func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    allowedOrigins := []string{
        "https://myapp.com",
        "https://www.myapp.com",
    }
    for _, allowed := range allowedOrigins {
        if origin == allowed {
            return true
        }
    }
    return false
}
```

### Message Size Limits

Default max message size: **512KB**

```go
const maxMessageSize = 512 * 1024 // Adjust as needed
connection.conn.SetReadLimit(maxMessageSize)
```

---

## Best Practices

### ✅ DO:

```javascript
// 1. Validate token before connecting
const token = localStorage.getItem('authToken')
if (!token) {
    // Redirect to login
}
wsStore.connect(token)

// 2. Handle reconnection
wsStore.onDisconnect(() => {
    setTimeout(() => wsStore.connect(token), 3000)
})

// 3. Subscribe to message types
wsStore.onMessage('notification', handleNotification)
wsStore.onMessage('chat', handleChat)

// 4. Cleanup on unmount (Svelte)
onDestroy(() => {
    wsStore.disconnect()
})
```

```go
// 5. Inject broadcaster into services
type OrderService struct {
    repo        ports.OrderRepository
    broadcaster *websocket.Broadcaster
    log         *zap.Logger
}

// 6. Send notifications from domain operations
func (s *OrderService) CompleteOrder(ctx context.Context, orderID string) error {
    if err := s.repo.UpdateStatus(ctx, orderID, "completed"); err != nil {
        return err
    }
    
    order, _ := s.repo.FindByID(ctx, orderID)
    s.broadcaster.NotifyUser(order.UserID, "order_completed", order)
    return nil
}
```

### ❌ DON'T:

```javascript
// ❌ Don't ignore disconnections
ws.onclose = () => { /* nothing */ }

// ❌ Don't send without validation
ws.send(untrustedData)

// ❌ Don't hold stale tokens
// Always refresh token before expiry
```

```go
// ❌ Don't broadcast sensitive data
broadcaster.Broadcast("user_data", user) // Password exposed!

// ❌ Don't ignore errors
broadcaster.SendToUser(userID, "msg", data) // No error check

// ❌ Don't block the Hub
func (h *Hub) Run() {
    select {
    case slowOperation := <-c: // WRONG! Can deadlock
        verySlowDatabaseCall()
    }
}
```

---

## Hub API Reference

### Broadcaster Methods

```go
// Inject into services via DI
type Broadcaster struct { hub *Hub }

// Broadcast to all clients
broadcaster.Broadcast(eventType, payload) error

// Send to specific user
broadcaster.NotifyUser(userID, eventType, payload) error

// Send to multiple users
broadcaster.NotifyMultipleUsers(userIDs, eventType, payload) error

// Send to room
broadcaster.BroadcastToRoom(roomID, eventType, payload) error

// Send error
broadcaster.SendError(userID, code, message) error

// Send status
broadcaster.SendStatus(message) error

// Query connection counts
broadcaster.ClientCount() int
broadcaster.UserConnCount(userID) int
broadcaster.RoomConnCount(roomID) int
```

---

## Troubleshooting

### Issue: WebSocket connection refused

**Causes:**
- Invalid token
- Server not running
- CORS origin mismatch (browser)

**Solution:**
```javascript
ws.onerror = (event) => {
    console.error('Connection failed:', event)
    // Token valid?
    // Server running?
}
```

### Issue: Messages not received

**Cause:** Client not listening on correct message type

**Solution:**
```javascript
// Make sure you subscribe to the right type
wsStore.onMessage('notification', handleIt)  // ✅
wsStore.onMessage('notification-wrong', ...) // ❌
```

### Issue: Memory leak / connections not cleaned up

**Cause:** readPump/writePump goroutines not exiting

**Solution:** Ensure connections are properly closed:
```go
defer func() {
    c.hub.unregister <- c
    c.conn.Close()
}()
```

### Issue: High CPU with many connections

**Cause:** Too many goroutines (readPump + writePump per connection)

**Solution:**
- Implement connection pooling
- Use Kubernetes horizontal scaling
- Monitor with `broadcaster.ClientCount()`

---

## Examples

### Full Chat Application

**Backend (Service):**
```go
type ChatService struct {
    repo        ports.ChatRepository
    broadcaster *websocket.Broadcaster
}

func (s *ChatService) SendChatMessage(ctx context.Context, roomID, userID, content string) error {
    msg := &domain.ChatMessage{
        RoomID: roomID,
        UserID: userID,
        Content: content,
        CreatedAt: time.Now(),
    }
    if err := s.repo.Create(ctx, msg); err != nil {
        return err
    }

    // Broadcast to all in room
    s.broadcaster.BroadcastToRoom(roomID, "chat", gin.H{
        "sender_id": userID,
        "content": content,
        "timestamp": msg.CreatedAt.Unix(),
    })

    return nil
}
```

**Frontend (Svelte):**
```svelte
<script>
    import { wsStore } from '$lib/stores/websocket'
    
    let messages = []
    let input = ''
    
    onMount(() => {
        wsStore.connect(token)
        wsStore.onMessage('chat', (msg) => {
            messages = [...messages, msg.payload]
        })
    })
    
    function send() {
        wsStore.send({
            type: 'chat',
            payload: { content: input }
        })
        input = ''
    }
</script>

<div>
    {#each messages as msg}
        <p>{msg.sender_id}: {msg.content}</p>
    {/each}
    <input bind:value={input} />
    <button on:click={send}>Send</button>
</div>
```

---

**Ready to add real-time features to your Kodia app! 🚀**
