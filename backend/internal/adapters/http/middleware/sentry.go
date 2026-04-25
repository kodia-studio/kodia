package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/getsentry/sentry-go"
)

// SentryTracing creates a Sentry transaction for each HTTP request, enabling performance monitoring per endpoint.
// Clones the Sentry hub per-request for goroutine-safe span tracking.
func SentryTracing(dsn string) gin.HandlerFunc {
	// If Sentry is not configured, return a no-op middleware
	if dsn == "" {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		// Clone the hub for request isolation (thread-safe)
		hub := sentry.CurrentHub().Clone()
		ctx := sentry.SetHubOnContext(c.Request.Context(), hub)
		c.Request = c.Request.WithContext(ctx)

		// Create a transaction for this HTTP request
		transactionName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		options := []sentry.SpanOption{
			sentry.WithOpName("http.server"),
			sentry.ContinueFromRequest(c.Request),
			sentry.WithTransactionName(transactionName),
		}

		transaction := sentry.StartTransaction(ctx, transactionName, options...)
		defer transaction.Finish()

		// Set request metadata on the transaction
		transaction.SetData("http.request.method", c.Request.Method)
		transaction.SetData("http.request.url", c.Request.URL.String())
		transaction.SetTag("http.route", c.FullPath())

		// Update context for handlers
		c.Request = c.Request.WithContext(transaction.Context())

		// Process the request
		c.Next()

		// Set response metadata and finish the transaction
		transaction.Status = httpStatusToSentryStatus(c.Writer.Status())
		transaction.SetData("http.response.status_code", c.Writer.Status())
	}
}

// httpStatusToSentryStatus maps HTTP status codes to Sentry span statuses.
func httpStatusToSentryStatus(code int) sentry.SpanStatus {
	switch {
	case code >= 200 && code < 300:
		return sentry.SpanStatusOK
	case code >= 300 && code < 400:
		return sentry.SpanStatusOK
	case code == 400:
		return sentry.SpanStatusInvalidArgument
	case code == 401:
		return sentry.SpanStatusUnauthenticated
	case code == 403:
		return sentry.SpanStatusPermissionDenied
	case code == 404:
		return sentry.SpanStatusNotFound
	case code >= 400 && code < 500:
		return sentry.SpanStatusFailedPrecondition
	case code >= 500:
		return sentry.SpanStatusInternalError
	default:
		return sentry.SpanStatusUnknown
	}
}
