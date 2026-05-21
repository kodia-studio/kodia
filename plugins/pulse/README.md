# Kodia Pulse Plugin

Real-time system monitoring and observability for Kodia Framework. Pulse provides live streaming of system metrics, performance data, and application logs.

## Features

- **Real-time System Metrics**: CPU, memory, disk usage, and goroutine count
- **Live Log Streaming**: Real-time warning and error logs from your application
- **WebSocket Integration**: Efficient bidirectional communication
- **Admin-Only Access**: Secured with JWT authentication and admin role requirement
- **Zero Configuration**: Works out of the box with minimal setup

## Installation

Add the Pulse plugin to your Kodia app initialization:

```go
// cmd/server/main.go
import "github.com/kodia-studio/pulse"

app.RegisterProviders(
    // ... other providers
    pulse.NewServiceProvider(),
)
```

## Usage

### Backend

The plugin automatically registers a WebSocket endpoint at `/api/pulse/stream` that requires admin authentication.

### Frontend

Access the Pulse dashboard at `http://localhost:5173/pulse` when running your Kodia application.

The dashboard displays:
- System vitals (CPU, RAM, Disk, Goroutines)
- CPU usage chart (last 30 samples)
- Memory usage chart (last 30 samples)
- Real-time system logs

## Architecture

- **Manager**: Orchestrates real-time telemetry broadcasting to connected clients
- **Handler**: WebSocket handler for client connections
- **Core**: Zap logging core that captures warn/error logs for streaming
- **ServiceProvider**: Plugin registration and lifecycle management

## Message Format

Messages are sent as JSON with the following format:

```json
{
  "type": "stats" | "log",
  "timestamp": "2026-05-19T10:30:00Z",
  "data": {
    // stats data or log data
  }
}
```

## Security

- Requires valid JWT token in query parameter or header
- Restricted to users with admin role
- WebSocket connections are authenticated and authorized
