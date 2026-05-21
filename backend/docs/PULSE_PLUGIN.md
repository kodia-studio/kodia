# Kodia Pulse Plugin 📊

Real-time system monitoring and observability for Kodia applications. Pulse provides live streaming of system metrics, performance data, and application logs through a professional admin dashboard.

## Overview

**Pulse** is an optional plugin that adds comprehensive monitoring capabilities to your Kodia application without any external dependencies. It streams live system vitals, memory usage, CPU metrics, and application logs directly to an admin dashboard.

### Key Features

- **Real-time Metrics**: CPU usage, memory, disk space, goroutine count
- **Live Charts**: Interactive visualizations of CPU and memory history
- **Log Streaming**: Real-time warning and error logs from your application
- **WebSocket Integration**: Efficient bidirectional communication
- **Zero Configuration**: Works out of the box with sensible defaults
- **Admin-Protected**: Secured with JWT authentication and role-based access

---

## Installation

### 1. Backend Setup

Add the Pulse plugin to your provider registration in `cmd/server/main.go`:

```go
import (
    // ... other imports
    "github.com/kodia-studio/pulse"
)

func main() {
    // ... initialization code ...
    
    err = app.RegisterProviders(
        providers.NewDatabaseProvider(),
        providers.NewInfraProvider(),
        providers.NewObservabilityProvider(),
        // ... other providers ...
        pulse.NewServiceProvider(), // Add Pulse plugin
    )
    
    // ... rest of boot sequence ...
}
```

### 2. Frontend Setup

Copy the Pulse dashboard page to your frontend routes:

```bash
cp plugins/pulse/frontend/routes/+page.svelte src/routes/(admin)/pulse/+page.svelte
```

Or if you're using a different route structure, place it at your admin dashboard path.

---

## Usage

Once enabled, access the Pulse dashboard at:

```
http://localhost:5173/admin/pulse
```

### Dashboard Components

#### System Vitals
Four key metrics displayed in a grid:
- **CPU Usage**: Current CPU utilization percentage
- **RAM Used**: Memory usage as percentage of total
- **Disk Space**: Disk usage across the system
- **Goroutines**: Number of active Go routines

#### Performance Charts
- **CPU Usage Chart**: Last 30 samples of CPU history
- **Memory Usage Chart**: Last 30 samples of memory history

Charts update in real-time as data arrives via WebSocket.

#### System Logs
A terminal-style log viewer displaying:
- **Timestamp**: When the log was generated
- **Level**: WARNING or ERROR
- **Message**: The log message content

The log stream is limited to 100 most recent entries, with new logs appearing at the top.

---

## Architecture

### Components

**Handler** (`handler.go`)
- Manages WebSocket connections
- Handles the `/api/pulse/stream` endpoint
- Authenticates and authorizes connections
- Manages message routing

**Manager** (`pulse_manager.go`)
- Orchestrates telemetry collection
- Gathers system stats every 2 seconds
- Broadcasts data to connected clients
- Manages client registration/unregistration
- Handles log streaming

**Core** (`pulse_core.go`)
- Integrates with Zap logger
- Intercepts warning and error logs
- Pipes logs to the manager for streaming

**Service Provider** (`service_provider.go`)
- Plugin lifecycle management
- Registers the manager in the app container
- Sets up the logger integration
- Registers HTTP routes

### Data Flow

```
┌─────────────────────┐
│   System Health     │
│   (CPU, Memory)     │
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  Pulse Manager      │
│  (Collect & Broadcast)
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  WebSocket Server   │
│  (/api/pulse/stream)│
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  Connected Clients  │
│  (Dashboard)        │
└─────────────────────┘
```

---

## Configuration

Pulse works with sensible defaults and requires no configuration:

| Setting | Default | Description |
|---------|---------|-------------|
| Update Interval | 2 seconds | How often stats are gathered |
| Log Buffer | 100 entries | Max logs kept in memory |
| Client Buffer | 256 messages | WebSocket message queue per client |
| Log Level | WARN | Minimum log level to capture |

To customize these, modify the values in `pulse_manager.go`:

```go
type Manager struct {
    // Adjust interval in Run()
    ticker := time.NewTicker(2 * time.Second)
    
    // Adjust log buffer size
    logs: make(chan LogData, 100),
}
```

---

## Message Format

Messages sent via WebSocket follow this JSON structure:

```json
{
  "type": "stats" | "log",
  "timestamp": "2026-05-19T10:30:00.123Z",
  "data": {
    // For stats:
    "cpu_usage_percent": 45.2,
    "memory_usage_percent": 62.1,
    "disk_usage_percent": 73.5,
    "goroutines": 128
    
    // For logs:
    "level": "error",
    "message": "Database connection timeout",
    "module": "database"
  }
}
```

---

## Security

### Authentication
- Requires valid JWT token in WebSocket connection
- Token can be passed as query parameter: `?token=<jwt_token>`

### Authorization
- Only users with the `admin` role can access Pulse
- Unauthorized access attempts are rejected

### Connection Security
- WebSocket connections are protected by auth middleware
- Each client connection is independent
- Connections are automatically cleaned up on disconnect

### In Production
```go
// Update WebSocket upgrader settings for production
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // Implement strict origin checking
        return r.Header.Get("Origin") == "https://yourdomain.com"
    },
}
```

---

## Performance Considerations

### Memory Usage
- Pulse maintains a small in-memory log buffer (100 entries by default)
- Client connections use ~256 messages per client queue
- Low overhead even with many concurrent connections

### CPU Impact
- Stats collection runs every 2 seconds
- Minimal CPU impact from telemetry gathering
- WebSocket handling is efficient and non-blocking

### Network
- Binary efficient JSON encoding
- Compresses well with gzip middleware
- Only admin users consume bandwidth (admin role)

---

## Troubleshooting

### Dashboard not loading
- Verify the frontend page is placed at `src/routes/(admin)/pulse/+page.svelte`
- Check that auth token is valid and user has admin role
- Check browser console for WebSocket errors

### WebSocket connection fails
- Verify the plugin is registered before Boot() is called
- Check that `/api/pulse/stream` endpoint is accessible
- Ensure JWT middleware is configured correctly
- Check firewall rules allow WebSocket connections

### No data appearing
- Check that the application is running
- Verify no errors in the server logs
- Confirm the auth token is not expired
- Check that user has admin role

### High memory usage
- Consider reducing log buffer size if logging excessively
- Monitor connected client count
- Check for disconnected clients not being cleaned up

---

## Example: Custom Integration

To extend Pulse with custom metrics:

```go
// In your code, access the manager
manager := kodia.MustResolve[*pulse.Manager](app, "pulse_manager")

// Log custom events (picked up by Pulse)
manager.Log("info", "Custom event occurred")

// Or send to dashboard directly via hooks
manager.send("custom", customData)
```

---

## Future Enhancements

Potential features for future versions:
- Custom metrics registration
- Configurable alert thresholds
- Historical data export
- Performance profiling integration
- Custom dashboard widgets

---

## Support

For issues or questions about the Pulse plugin:
1. Check the troubleshooting section above
2. Review application logs for errors
3. Verify plugin configuration
4. Open an issue on the Kodia Framework GitHub repository

---

© 2026 Kodia Studio. "Build like a user, code like a pro."
