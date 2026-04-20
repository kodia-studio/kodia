# Observability & Monitoring

Kodia provides a state-of-the-art observability stack that gives you deep visibility into your application's performance, health, and errors.

## Overview

The Kodia observability stack is built on four pillars:

1.  **Distributed Tracing**: Powered by OpenTelemetry (OTEL).
2.  **Metrics Collection**: Powered by Prometheus.
3.  **Real-time Health Monitoring**: Built-in system stats gatherer.
4.  **Error Tracking**: Native Sentry integration.

---

## 1. Distributed Tracing

Tracing allows you to see the lifecycle of a request as it moves through your system. Kodia uses OpenTelemetry to automatically instrument every HTTP request.

### Configuration
Enable tracing in your `config.yaml` or via environment variables:

```yaml
observability:
  tracing_enabled: true
  service_name: "my-kodia-api"
  otlp_endpoint: "localhost:4317" # Jaeger, Datadog, or New Relic OTLP collector
  sampling_rate: 1.0 # 100% of traces
```

---

## 2. Metrics (Prometheus)

Kodia collects real-time metrics for your application, including request counts and latency histograms.

### Scraping Metrics
By default, Kodia starts a dedicated metrics server on port `9090`. You can scrape it using Prometheus at the `/metrics` endpoint.

---

## 3. Health Checks

Kodia monitors the physical health of its environment, including CPU, Memory, and Disk usage.

### CLI Health Command
For a quick check from your terminal, use the Kodia CLI:
```bash
kodia health
```

---

## 4. Error Tracking (Sentry)

If an unexpected crash occurs, Kodia automatically captures the stack trace and reports it to Sentry.

### Configuration
```yaml
observability:
  sentry_dsn: "https://your-dsn-here@sentry.io/123"
```

---

## 5. Performance Profiling (pprof)

In development mode, Kodia exposes profiling endpoints at `/debug/pprof`.

```bash
go tool pprof http://localhost:8080/debug/pprof/profile
```
