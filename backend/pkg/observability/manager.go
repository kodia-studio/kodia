package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// TODO: Re-enable OTLP tracing when the google.golang.org/genproto split-module
// conflict with grpc-gateway is resolved upstream. Tracked issue:
// https://github.com/open-telemetry/opentelemetry-go/issues/XXXX
// Temporarily disabled to unblock the build. Sentry + Prometheus remain active.

// Manager handles the initialization and shutdown of observability tools.
type Manager struct {
	cfg *config.Config
	log *zap.Logger
}

// NewManager creates a new Observability Manager.
func NewManager(cfg *config.Config, log *zap.Logger) *Manager {
	return &Manager{
		cfg: cfg,
		log: log,
	}
}

// Init starts the tracing, metrics, and sentry clients.
func (m *Manager) Init(ctx context.Context) error {
	// 1. Initialize Sentry
	if m.cfg.Observability.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              m.cfg.Observability.SentryDSN,
			Environment:      m.cfg.App.Env,
			TracesSampleRate: m.cfg.Observability.SamplingRate,
		})
		if err != nil {
			m.log.Error("Sentry initialization failed", zap.Error(err))
		} else {
			m.log.Info("Sentry initialized successfully")
		}
	}

	// 2. OpenTelemetry (Tracing) — temporarily disabled due to genproto split-module conflict.
	// TODO: Re-enable when upstream otel/grpc-gateway conflict is resolved.
	if m.cfg.Observability.TracingEnabled {
		m.log.Warn("OpenTelemetry OTLP Tracing is temporarily disabled due to module conflict",
			zap.String("endpoint", m.cfg.Observability.OTLPEndpoint))
	}

	// 3. Initialize Prometheus (Metrics)
	if m.cfg.Observability.MetricsEnabled {
		go m.startMetricsServer()
		m.log.Info("Prometheus Metrics server starting", zap.Int("port", m.cfg.Observability.PrometheusPort))
	}

	return nil
}

func (m *Manager) startMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	
	addr := fmt.Sprintf(":%d", m.cfg.Observability.PrometheusPort)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		m.log.Error("Prometheus server failed", zap.Error(err))
	}
}

// Shutdown gracefully shuts down the observability stack.
func (m *Manager) Shutdown(ctx context.Context) {
	// TracerProvider shutdown disabled — tracing temporarily suspended.
	_ = ctx
	sentry.Flush(2 * time.Second)
}
