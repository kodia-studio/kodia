# Observability & DevOps

Kodia provides a production-grade suite for monitoring, diagnostics, and reliable deployments. It includes health checks, distributed tracing, metrics, structured audit logging, and graceful shutdown out of the box.

---

## 1. Health Checks (Liveness & Readiness)

Kodia follows Kubernetes standards for health checks with modular checkers for your dependencies.

### Endpoints
- `GET /api/v1/health/live` — Returns 200 if the process is running.
- `GET /api/v1/health/ready` — Checks database, redis, and other critical dependencies.

### Configuration
Health checks are managed via `ObservabilityProvider`. It automatically detects if Database or Redis is registered and adds the corresponding checkers.

### Adding Custom Checkers
You can implement the `health.Checker` interface and add it to your custom provider:

```go
type MyServiceChecker struct{}
func (c *MyServiceChecker) Name() string { return "my_service" }
func (c *MyServiceChecker) Check(ctx context.Context) error {
    // Check if your service is healthy
    return nil
}
```

---

## 2. Distributed Tracing (OpenTelemetry)

Kodia integrates **OpenTelemetry (OTEL)** to trace requests as they move through your application services.

### How it works
- **Middleware**: Automatically creates a trace span for every incoming HTTP request.
- **Propagation**: Injects `X-Trace-ID` into response headers.
- **Exporting**: By default, Kodia uses the `stdout` exporter with pretty-print for high visibility during development.

### Production Setup
To export to an external collector (like Jaeger, Honeycomb, or Tempo), swap the exporter in `pkg/observability/manager.go` to `otlptracehttp`.

### Environment Variables
```env
APP_OBSERVABILITY_TRACING_ENABLED=true
APP_OBSERVABILITY_OTLP_ENDPOINT=localhost:4318
APP_OBSERVABILITY_SAMPLING_RATE=1.0
```

---

## 3. Metrics (Prometheus)

Kodia exposes a Prometheus-compatible metrics endpoint for scraping.

### Endpoint
- `GET :9090/metrics` (Default port)

### Available Metrics
- **Go Runtime**: Memory usage, goroutine count, GC stats.
- **HTTP**: Request count, duration, status codes (via middleware).
- **Custom**: You can register your own counters and gauges via `prometheus` package.

---

## 4. Structured Audit Log

Kodia provides an enterprise-grade audit logging system that records every important action.

### Multi-Sink Support
The `AuditManager` broadcasts logs to multiple destinations:
1. **Database**: Saved to the `audit_logs` table (GORM).
2. **Structured Log**: Output to system logs as JSON (Zap).

### Usage
```go
auditManager := app.MustGet("audit").(*audit.Manager)

auditManager.LogAction(
    userID, 
    userEmail, 
    "Order #123", 
    audit.ActionUpdate, 
    audit.StatusSuccess, 
    "Changed status to shipped", 
    ipAddress, 
    userAgent,
)
```

---

## 5. Graceful Shutdown

Zero-downtime deployment is achieved by handling OS signals and orchestrating resource cleanup.

### How it works
When the application receives `SIGINT` or `SIGTERM`:
1. The HTTP server stops accepting new connections but finishes active ones (30s timeout).
2. **Cleanup Tasks** are executed in reverse order of registration:
    - Closing Database connections.
    - Closing Redis connections.
    - Flushing Sentry and OpenTelemetry tracers.

### Registering Cleanup Tasks
You can add your own cleanup logic in any provider:
```go
app.RegisterCleanupTask(func(ctx context.Context) error {
    app.Log.Info("Cleaning up my resource...")
    return myResource.Close()
})
```

---

## 6. Provider Setup

Ensure `ObservabilityProvider` is registered in your `main.go`:

```go
app.RegisterProviders(
    providers.NewDatabaseProvider(),
    providers.NewInfraProvider(),
    providers.NewObservabilityProvider(), // Required
    providers.NewHttpProvider(),
    // ...
)
```
