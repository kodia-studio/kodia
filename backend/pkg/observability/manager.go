package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)


// Manager handles the initialization and shutdown of observability tools.
type Manager struct {
	cfg            *config.Config
	log            *zap.Logger
	tracerProvider *trace.TracerProvider
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

	// 2. OpenTelemetry (Tracing) — Using HTTP OTLP to avoid gRPC dependency conflicts
	if m.cfg.Observability.TracingEnabled {
		if err := m.initTracer(ctx); err != nil {
			m.log.Error("OpenTelemetry initialization failed", zap.Error(err))
		} else {
			m.log.Info("OpenTelemetry Tracing initialized", zap.String("endpoint", m.cfg.Observability.OTLPEndpoint))
		}
	}

	// 3. Initialize Prometheus (Metrics)
	if m.cfg.Observability.MetricsEnabled {
		go m.startMetricsServer()
		m.log.Info("Prometheus Metrics server starting", zap.Int("port", m.cfg.Observability.PrometheusPort))
	}

	return nil
}

func (m *Manager) initTracer(ctx context.Context) error {
	// Using stdouttrace for zero-conflict, high-visibility debugging.
	// In production, this can be swapped for otlptracehttp.
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create stdout trace exporter: %w", err)
	}


	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.cfg.Observability.ServiceName),
			semconv.DeploymentEnvironmentKey.String(m.cfg.App.Env),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(trace.TraceIDRatioBased(m.cfg.Observability.SamplingRate)),
	)
	otel.SetTracerProvider(tp)
	m.tracerProvider = tp

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
	if m.tracerProvider != nil {
		if err := m.tracerProvider.Shutdown(ctx); err != nil {
			m.log.Error("OpenTelemetry TracerProvider shutdown failed", zap.Error(err))
		}
	}
	sentry.Flush(2 * time.Second)
}
