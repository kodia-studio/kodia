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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

// Manager handles the initialization and shutdown of observability tools.
type Manager struct {
	cfg            *config.Config
	log            *zap.Logger
	tracerProvider *sdktrace.TracerProvider
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

	// 2. Initialize OpenTelemetry (Tracing)
	if m.cfg.Observability.TracingEnabled {
		if err := m.initTracing(ctx); err != nil {
			return fmt.Errorf("tracing init: %w", err)
		}
		m.log.Info("OpenTelemetry Tracing initialized", zap.String("endpoint", m.cfg.Observability.OTLPEndpoint))
	}

	// 3. Initialize Prometheus (Metrics)
	if m.cfg.Observability.MetricsEnabled {
		go m.startMetricsServer()
		m.log.Info("Prometheus Metrics server starting", zap.Int("port", m.cfg.Observability.PrometheusPort))
	}

	return nil
}

func (m *Manager) initTracing(ctx context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(m.cfg.Observability.OTLPEndpoint),
		otlptracehttp.WithInsecure(), // Adjust if using TLS in prod
	)
	if err != nil {
		return err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.cfg.Observability.ServiceName),
			semconv.DeploymentEnvironmentKey.String(m.cfg.App.Env),
		),
	)
	if err != nil {
		return err
	}

	m.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(m.cfg.Observability.SamplingRate)),
	)

	otel.SetTracerProvider(m.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

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
			m.log.Error("TracerProvider shutdown failed", zap.Error(err))
		}
	}
	sentry.Flush(2 * time.Second)
}
