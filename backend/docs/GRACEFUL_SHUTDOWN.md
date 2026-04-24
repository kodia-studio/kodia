# Graceful Shutdown

Kodia Framework implements production-grade graceful shutdown that safely stops the application while allowing in-flight requests to complete.

## Overview

Graceful shutdown ensures:
- ✅ In-flight requests complete before shutdown
- ✅ No new connections accepted after shutdown signal
- ✅ Database connections properly closed
- ✅ Redis connections released
- ✅ Background tasks completed
- ✅ Configurable timeout (default: 30 seconds)

## How It Works

### 1. Signal Handling

The application listens for termination signals:

```go
// Signals that trigger graceful shutdown
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
```

**Supported Signals:**
- `SIGINT` (Ctrl+C) — User interrupt
- `SIGTERM` (kill) — Kubernetes, Docker, systemd
- `SIGHUP` (kill -1) — Hot reload, configuration changes

### 2. Shutdown Sequence

When a signal is received:

```
1. Stop accepting new connections
2. Drain existing connections (timeout: configurable)
3. Close HTTP server
4. Run cleanup tasks (in reverse order)
5. Exit
```

### 3. Shutdown Logging

Each phase logs detailed information:

```
2026-04-24T10:30:45Z INFO  Signal received signal=SIGTERM
2026-04-24T10:30:45Z INFO  Draining connections... timeout=30s
2026-04-24T10:30:46Z INFO  Running cleanup tasks... count=5
2026-04-24T10:30:47Z INFO  Shutdown complete
```

## Configuration

### Timeout Configuration

**Default:** 30 seconds

**Via Environment:**
```bash
APP_SHUTDOWN_TIMEOUT_SECS=60  # Wait up to 1 minute
```

**In Code:**
```go
cfg, _ := config.Load()
cfg.App.ShutdownTimeoutSecs  // reads from env or default
```

**Example Scenarios:**

| Timeout | Use Case |
|---------|----------|
| 10s | Development, quick restarts |
| 30s | Standard (default) |
| 60s | Large database operations |
| 120s | Complex cleanup tasks |

### Setting Different Timeouts

```bash
# Development (quick)
APP_SHUTDOWN_TIMEOUT_SECS=10 go run ./cmd/server

# Production (generous)
APP_SHUTDOWN_TIMEOUT_SECS=60 ./server
```

## Implementation

### Source Code

**File:** `pkg/kodia/app.go`

```go
func (a *App) Run() error {
    // ... start server ...

    // Listen for shutdown signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    sig := <-quit

    a.Log.Info("Signal received", zap.String("signal", sig.String()))

    // Configure timeout from config (or use 30s default)
    timeout := time.Duration(a.Config.App.ShutdownTimeoutSecs) * time.Second
    if timeout == 0 {
        timeout = 30 * time.Second
    }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // Drain connections
    a.Log.Info("Draining connections...", zap.Duration("timeout", timeout))
    if err := server.Shutdown(ctx); err != nil {
        a.Log.Error("HTTP server forced shutdown", zap.Error(err))
    }

    // Run cleanup tasks
    a.Log.Info("Running cleanup tasks...", zap.Int("count", len(a.cleanupTasks)))
    for i := len(a.cleanupTasks) - 1; i >= 0; i-- {
        if err := a.cleanupTasks[i](ctx); err != nil {
            a.Log.Error("Cleanup task failed", zap.Error(err))
        }
    }

    a.Log.Info("Shutdown complete")
    return nil
}
```

### Cleanup Tasks

Providers register cleanup tasks during initialization:

```go
// In DatabaseProvider.Register()
app.RegisterCleanupTask(func(ctx context.Context) error {
    sqlDB, _ := db.DB()
    app.Log.Info("Closing database connection...")
    return sqlDB.Close()
})

// In RedisProvider.Register()
app.RegisterCleanupTask(func(ctx context.Context) error {
    app.Log.Info("Closing Redis connection...")
    return redisClient.Close()
})
```

**Cleanup Task Order:**
Tasks are executed in REVERSE registration order (LIFO):

```
1. Register: Database (first)
2. Register: Redis (second)
3. Register: Cache (third)

Shutdown order:
1. Cache cleanup (last registered, first cleaned up)
2. Redis cleanup
3. Database cleanup (first registered, last cleaned up)
```

This ensures dependencies are shut down in the correct order.

## Deployment Scenarios

### 1. Kubernetes Pod Termination

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: kodia-app
spec:
  terminationGracePeriodSeconds: 60  # Must be >= APP_SHUTDOWN_TIMEOUT_SECS
  containers:
  - name: app
    image: kodia:latest
    env:
    - name: APP_SHUTDOWN_TIMEOUT_SECS
      value: "45"
```

**Sequence:**
1. Kubectl sends SIGTERM
2. App has 45s to gracefully shutdown (per APP_SHUTDOWN_TIMEOUT_SECS)
3. App must complete shutdown within 60s (per terminationGracePeriodSeconds)
4. If not done by 60s, SIGKILL is sent

### 2. Docker Container Stop

```bash
# Docker sends SIGTERM, waits 10s, then SIGKILL
docker stop --time 60 kodia-container

# Run container with custom signal
docker run --stop-signal SIGTERM kodia:latest
```

**Best Practice:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
COPY --from=builder /app/server .
ENV APP_SHUTDOWN_TIMEOUT_SECS=30
EXPOSE 8080
CMD ["./server"]
```

### 3. Systemd Service

```ini
[Unit]
Description=Kodia API Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/kodia
ExecStop=/bin/kill -TERM $MAINPID
TimeoutStopSec=60
Restart=on-failure
RestartSec=10
Environment="APP_SHUTDOWN_TIMEOUT_SECS=45"

[Install]
WantedBy=multi-user.target
```

### 4. Manual Stop

```bash
# Graceful shutdown
kill -TERM $PID

# If graceful shutdown fails after timeout, force kill
kill -9 $PID
```

## What Happens During Shutdown

### Active Request Example

```
Time: T+0s
- App receives SIGTERM
- Logging: "Signal received signal=SIGTERM"

Time: T+0s-T+25s
- Request still processing in handler
- HTTP server stops accepting NEW connections
- Handler completes and returns response

Time: T+25s
- All in-flight requests finished
- Logging: "Draining connections... timeout=30s"
- HTTP server shutdown complete

Time: T+25s-T+28s
- Database connection closed
- Logging: "Running cleanup tasks... count=5"
- All cleanup tasks completed

Time: T+28s
- Logging: "Shutdown complete"
- Process exits with code 0
```

### Timeout Exceeded Example

```
Time: T+0s
- App receives SIGTERM
- Logging: "Signal received signal=SIGTERM"

Time: T+0s-T+30s
- Long-running database query still executing
- Timeout waiting for request to finish

Time: T+30s
- Timeout exceeded
- Logging: "HTTP server forced shutdown error=context deadline exceeded"
- Continue cleanup anyway

Time: T+30s-T+32s
- Database connection closed (may interrupt query)
- Logging: "Cleanup task failed error=context deadline exceeded"
- Process exits with code 1
```

## Monitoring Shutdown

### Metrics to Track

```bash
# Time taken for graceful shutdown
shutdown_duration_seconds

# Count of in-flight requests during shutdown
in_flight_requests_count

# Failed cleanup tasks during shutdown
failed_cleanup_tasks_count

# Forced shutdowns (timeout exceeded)
forced_shutdown_count
```

### Log Analysis

```bash
# Find all graceful shutdowns
grep "Signal received" app.log

# Find forced shutdowns (errors during shutdown)
grep "forced shutdown" app.log

# Find slow shutdowns (>20s)
grep "Shutdown complete" app.log | grep "T+[2-9][0-9]s"
```

## Best Practices

### ✅ DO:

- Set appropriate `APP_SHUTDOWN_TIMEOUT_SECS` based on workload
- Ensure Kubernetes `terminationGracePeriodSeconds >= shutdown_timeout_secs`
- Monitor logs for "Signal received" to track restarts
- Test graceful shutdown in staging before production
- Use SIGTERM for safe shutdowns (not SIGKILL)
- Log context information (request IDs) during shutdown
- Handle context timeouts in cleanup tasks

### ❌ DON'T:

- Use SIGKILL (9) for normal shutdown — this skips cleanup
- Set timeout to 0 or negative values
- Keep long-running tasks that exceed timeout
- Ignore "forced shutdown" errors in logs
- Change timeouts too frequently
- Skip cleanup task registration
- Run blocking I/O in signal handlers

## Troubleshooting

### Issue: Shutdown Takes Too Long

**Diagnosis:**
```bash
grep "Shutdown complete" app.log
# If time between "Signal received" and "Shutdown complete" > timeout
```

**Solutions:**
1. Increase `APP_SHUTDOWN_TIMEOUT_SECS`
2. Optimize cleanup tasks to be faster
3. Check for slow database queries
4. Review request handlers for long operations

### Issue: Shutdown Always Times Out

**Diagnosis:**
```bash
grep "forced shutdown" app.log
# If appears in every restart
```

**Solutions:**
1. Increase timeout significantly
2. Add context timeout checks in handlers
3. Use worker goroutines that respect context cancellation
4. Break long operations into smaller chunks

### Issue: Database Connection Not Closed

**Diagnosis:**
```bash
# Monitor connection count during shutdown
psql postgres -c "SELECT datname, count(*) FROM pg_stat_activity GROUP BY datname"
# Should decrease to 0 during shutdown
```

**Solutions:**
1. Ensure DatabaseProvider registers cleanup task
2. Check database pool configuration
3. Verify cleanup task executes without error

## Testing Graceful Shutdown

### Unit Test

```go
func TestGracefulShutdown(t *testing.T) {
    app := setupTestApp()

    // Start in goroutine
    go app.Run()
    time.Sleep(100 * time.Millisecond)  // Let it start

    // Send SIGTERM
    proc, _ := os.FindProcess(os.Getpid())
    proc.Signal(syscall.SIGTERM)

    // App should exit gracefully
    // (in real tests, use signal channels)
}
```

### Integration Test

```bash
# Start server
./server &
SERVER_PID=$!

# Wait for startup
sleep 1

# Send SIGTERM
kill -TERM $SERVER_PID

# Wait for shutdown
wait $SERVER_PID
EXIT_CODE=$?

# Verify clean shutdown
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ Graceful shutdown successful"
else
    echo "✗ Shutdown failed with code $EXIT_CODE"
fi
```

### Load Test During Shutdown

```bash
# Start server
./server &
SERVER_PID=$!

# Start continuous requests
while true; do
    curl http://localhost:8080/api/v1/health &
done &

# After 10s, shutdown
sleep 10
kill -TERM $SERVER_PID

# All pending requests should complete
# No errors in logs
```

## Deployment Checklist

Before deploying to production:

- [ ] Understand your workload's longest request duration
- [ ] Set `APP_SHUTDOWN_TIMEOUT_SECS >= max_request_duration + buffer`
- [ ] Configure orchestrator (K8s, Docker, systemd) shutdown times
- [ ] Test graceful shutdown in staging environment
- [ ] Set up monitoring/alerts for shutdown duration
- [ ] Document shutdown procedure for ops team
- [ ] Have plan B for forced shutdown if grace period exceeded
- [ ] Log all shutdown signals and durations
- [ ] Verify all cleanup tasks are registered
- [ ] Ensure cleanup tasks are idempotent

## References

- [Go graceful shutdown pattern](https://golang.org/doc/effective_go#channels)
- [Kubernetes termination grace period](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#termination-of-pods)
- [Docker container stop](https://docs.docker.com/engine/reference/commandline/stop/)
- [systemd service termination](https://www.freedesktop.org/software/systemd/man/systemd.service.html)
- [Signal handling in Go](https://golang.org/pkg/os/signal/)

## Conclusion

Graceful shutdown in Kodia ensures:
- 🛡️ **Data Integrity**: Database transactions complete safely
- 🛡️ **Request Reliability**: In-flight requests finish and respond
- 🛡️ **Zero Downtime**: Deployment restarts without errors
- 📊 **Observability**: Full logging of shutdown sequence
- ⚡ **Production Ready**: Orchestrator-aware timeout handling
