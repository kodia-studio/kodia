package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracing returns a Gin middleware that instruments requests with OpenTelemetry.
func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer("kodia-http")

	return func(c *gin.Context) {
		// Extract existing context from headers (for distributed tracing)
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start a new span
		opts := []trace.SpanStartOption{
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", c.FullPath()),
				attribute.String("http.remote_addr", c.ClientIP()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		}

		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// Inject span into request context
		c.Request = c.Request.WithContext(ctx)

		// Set trace ID in response header for debugging
		if span.SpanContext().HasTraceID() {
			c.Header("X-Trace-ID", span.SpanContext().TraceID().String())
		}

		c.Next()

		// Record status code and errors
		status := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", status))
		
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("http.error", c.Errors.String()))
			span.RecordError(fmt.Errorf("%s", c.Errors.String()))
		}
	}
}
