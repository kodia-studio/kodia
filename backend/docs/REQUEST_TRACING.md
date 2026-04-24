# Request Tracing with X-Request-ID

Kodia Framework provides built-in request tracing to help debug, monitor, and audit API requests end-to-end.

## Overview

Every request flowing through the Kodia API is assigned a unique identifier (`request_id`) that:
- ✅ Identifies the request uniquely across all logs
- ✅ Appears in all response envelopes
- ✅ Included in structured logs for debugging
- ✅ Can be passed from client to trace request chains

## How It Works

### 1. Request ID Generation

**Middleware:** `middleware/request_id.go`

The `RequestID()` middleware (registered first in the middleware chain):

```go
// Registered FIRST, before all other middleware
engine.Use(middleware.RequestID())
engine.Use(middleware.Recovery(r.log))
engine.Use(middleware.Logger(r.log))
```

**Process:**
1. Check for `X-Request-ID` header in incoming request
2. If present, use that value (allows client-provided tracing)
3. If not present, generate new UUID4
4. Store in gin context with key `"request_id"`
5. Set `X-Request-ID` response header

### 2. Request ID in Responses

**File:** `pkg/response/response.go`

All API responses include `request_id` field:

```json
{
  "success": true,
  "message": "User created successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

**Error responses also include request_id:**

```json
{
  "success": false,
  "message": "Validation failed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "errors": {
    "email": ["Email is required"]
  }
}
```

### 3. Request ID in Logs

**File:** `middleware/logger.go`

All structured logs include `request_id` field:

```json
{
  "timestamp": "2026-04-24T10:30:45.123Z",
  "level": "INFO",
  "message": "Request",
  "status": 201,
  "method": "POST",
  "path": "/api/v1/auth/register",
  "ip": "203.0.113.195",
  "latency": 125,
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "42"
}
```

## Usage Patterns

### 1. Client-Side Tracing

Clients can send custom request IDs for end-to-end tracing:

```bash
# Pass request ID from client
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "X-Request-ID: my-app-123-abc" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret"}'

# Response includes the same request ID
{
  "success": true,
  "request_id": "my-app-123-abc",
  "data": { "token": "..." }
}
```

### 2. Server-Side Tracing

When Kodia is called from another service:

```go
// Service A calls Service B with tracing
req, _ := http.NewRequest("POST", "http://serviceB:8080/api/users", body)
req.Header.Set("X-Request-ID", "parent-request-id-123")

// Service B automatically forwards this request ID
// All logs and responses include "parent-request-id-123"
```

### 3. Debugging Failed Requests

When a user reports a problem, they include the request ID:

```
User: "The login failed with error XYZ at 10:30 AM"
Support: "Please provide the request ID shown in the error"
User: "It's 550e8400-e29b-41d4-a716-446655440000"
```

Support team searches logs:

```bash
# Find all logs for this request
grep "request_id.*550e8400-e29b-41d4-a716-446655440000" app.log

# Output:
# 10:30:44 - Request received
# 10:30:44 - Database query executed
# 10:30:45 - Email validation failed
# 10:30:45 - Response sent with 400 Bad Request
```

## Architecture

### Request ID Flow

```
┌─────────────────────────────────────────────────────────────┐
│ Client Request                                              │
│ (with or without X-Request-ID header)                       │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ RequestID Middleware (FIRST)                                │
│ ├─ Check for X-Request-ID header                            │
│ ├─ Generate UUID4 if missing                                │
│ └─ Store in gin.Context                                     │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ Logger Middleware                                           │
│ └─ Extract request_id from context and add to logs          │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ Handler / Endpoint                                          │
│ (has access to request_id via gin.Context)                 │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ Response Envelope (pkg/response)                            │
│ └─ Inject request_id from context into response JSON        │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ Response Headers + Body                                     │
│ ├─ X-Request-ID header                                      │
│ └─ JSON body with request_id field                          │
└─────────────────────────────────────────────────────────────┘
```

## Implementation Details

### Middleware Order (Important!)

RequestID middleware MUST be first:

```go
// ✅ CORRECT - RequestID first
engine.Use(middleware.RequestID())        // 1st
engine.Use(middleware.Recovery(r.log))    // 2nd
engine.Use(middleware.Logger(r.log))      // 3rd
engine.Use(cors.New(corsConfig))          // 4th
// ... other middleware ...
```

If RequestID is not first, earlier middleware won't have access to request_id.

### Accessing Request ID in Handlers

```go
func (h *UserHandler) Create(c *gin.Context) {
    // Get request ID from context
    requestID := c.GetString(middleware.RequestIDKey)
    
    // Use in logging
    h.log.Info("Creating user", 
        zap.String("request_id", requestID),
        zap.String("email", req.Email),
    )
    
    // Pass to services if needed
    user, err := h.service.Create(c.Request.Context(), req, requestID)
    if err != nil {
        response.BadRequest(c, "Validation failed", err)
        return
    }
    
    response.Created(c, "User created", user)
}
```

### Automatic Injection in Responses

The `response` package automatically injects request_id:

```go
// In pkg/response/response.go
func send(c *gin.Context, status int, r Response) {
    // Automatically inject request_id from context
    if reqID := c.GetString("request_id"); reqID != "" {
        r.RequestID = reqID
    }
    c.JSON(status, r)
}
```

No manual work needed — just use standard response functions:

```go
response.Created(c, "User created", user)  // request_id auto-included!
response.BadRequest(c, "Invalid email", errors)
response.InternalServerError(c, "Database error")
```

## Configuration

Request ID can be controlled via environment:

```bash
# Disable request ID generation (not recommended)
# Currently no env var — always enabled for production readiness

# Custom header name (future feature)
# Currently hardcoded to X-Request-ID
```

## Monitoring & Analytics

### Search Request Trails

**ELK Stack:**
```
# Find all requests from user
GET /logs/_search
{
  "query": {
    "match": { "user_id": "42" }
  },
  "sort": [{ "timestamp": "asc" }]
}
```

**Datadog:**
```
# Group by request_id
fields @request_id, @message
| stats count() as num_events by @request_id
| sort num_events desc
```

### Correlate With Client Sessions

```bash
# If client sends request ID like: "session-123-req-1"
# You can trace which requests belong to which session
filter(request_id startswith "session-123")
```

## Best Practices

### ✅ DO:

- Always include request ID in error messages shown to users
- Log request_id in all audit trails
- Use request IDs for distributed tracing between services
- Include request_id in support/bug reports
- Monitor request_id length and format (should be valid UUID4)
- Preserve request_id across service boundaries

### ❌ DON'T:

- Don't expose internal request IDs in public documentation
- Don't use request IDs for security purposes (not cryptographic)
- Don't rely on request_id for duplicate detection (same UUID could happen)
- Don't send request_id across unencrypted connections
- Don't filter/remove request_id from responses

## Performance Impact

- **Middleware overhead**: <1ms per request
- **Memory impact**: ~36 bytes per request (UUID string)
- **Log storage impact**: ~50 bytes per log line (request_id field)

Negligible impact on performance.

## Security Considerations

### Request ID Leakage

Request IDs are **not sensitive** and can be logged/displayed:
- ✅ Safe to show in error messages
- ✅ Safe to include in emails
- ✅ Safe to log in access logs
- ⚠️ Don't use as authentication token

### Attack Surface

Request IDs don't increase attack surface:
- Generated with strong randomness (UUID4)
- Only used for logging/tracing
- Cannot be exploited for CSRF, XSS, or injection
- Cannot be used for privilege escalation

## Testing

### Test Request ID Propagation

```go
func TestRequestIDPropagation(t *testing.T) {
    app := setupTestApp()
    
    // Test 1: Auto-generated request ID
    w := testRequest(app, "GET", "/api/v1/health", nil)
    
    requestID := w.Header().Get("X-Request-ID")
    assert.NotEmpty(t, requestID)
    assert.Regexp(t, `^[0-9a-f-]{36}$`, requestID)  // UUID format
    
    // Test 2: Client-provided request ID
    req := httptest.NewRequest("GET", "/api/v1/health", nil)
    req.Header.Set("X-Request-ID", "custom-123")
    w = httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, "custom-123", w.Header().Get("X-Request-ID"))
}
```

## Compliance

### Standards & Regulations

- ✅ **HIPAA**: Request IDs help with audit trails
- ✅ **PCI-DSS**: Supports comprehensive logging requirement
- ✅ **SOC 2 Type II**: Evidence of request tracing capability
- ✅ **GDPR**: Links user actions to requests for data access requests

## References

- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
- [OpenTelemetry Specification](https://opentelemetry.io/docs/reference/specification/)
- [UUID RFC 4122](https://tools.ietf.org/html/rfc4122)
- [HTTP Header X-Request-ID](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers)

## Conclusion

Request tracing with X-Request-ID provides:
- 🔍 Complete audit trail for every request
- 🐛 Easy debugging with full request lifecycle
- 📊 Better monitoring and analytics
- 🔐 Compliance with regulatory requirements
- ⚡ Minimal performance overhead
