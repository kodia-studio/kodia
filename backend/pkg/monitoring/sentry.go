// Package monitoring provides developer-friendly wrappers for Sentry error monitoring.
package monitoring

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// CaptureError reports an error to Sentry with optional context tags.
// Returns the event ID, or nil if Sentry is not initialized.
func CaptureError(ctx context.Context, err error, tags map[string]string) *sentry.EventID {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil || hub.Client() == nil {
		return nil
	}

	if tags != nil {
		scope := hub.Scope()
		for key, value := range tags {
			scope.SetTag(key, value)
		}
	}

	return hub.CaptureException(err)
}

// CaptureMessage sends an informational or warning message to Sentry.
// Level can be: sentry.LevelDebug, sentry.LevelInfo, sentry.LevelWarning, sentry.LevelError, sentry.LevelFatal
// Returns the event ID, or nil if Sentry is not initialized.
func CaptureMessage(ctx context.Context, msg string, level sentry.Level) *sentry.EventID {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil || hub.Client() == nil {
		return nil
	}

	scope := hub.Scope()
	scope.SetLevel(level)
	return hub.CaptureMessage(msg)
}

// SentryUser represents user information to attach to Sentry events.
type SentryUser struct {
	ID       string
	Email    string
	Username string
	IPAddress string
}

// SetUser attaches user context to the active Sentry scope.
// Call this in your authentication middleware after the user is verified.
func SetUser(ctx context.Context, user SentryUser) {
	if !sentry.HasHubOnContext(ctx) {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	if hub.Client() == nil {
		return
	}

	hub.Scope().SetUser(sentry.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		IPAddress: user.IPAddress,
	})
}

// AddBreadcrumb records a breadcrumb event on the Sentry scope.
// Breadcrumbs are useful for tracking important events leading up to an error.
func AddBreadcrumb(ctx context.Context, category, message string, level sentry.Level, data map[string]interface{}) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil || hub.Client() == nil {
		return
	}

	hub.Scope().AddBreadcrumb(&sentry.Breadcrumb{
		Category: category,
		Message:  message,
		Level:    level,
		Data:     data,
	}, -1)
}

// WithSpan wraps a function in a Sentry span for performance profiling.
// The span operation and description help organize performance metrics in Sentry.
// Returns the error from fn, if any.
func WithSpan(ctx context.Context, op, description string, fn func(context.Context) error) error {
	if !sentry.HasHubOnContext(ctx) {
		return fn(ctx)
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil || hub.Client() == nil {
		return fn(ctx)
	}

	span := sentry.StartSpan(ctx, op, sentry.WithDescription(description))
	defer span.Finish()

	return fn(span.Context())
}

// IsInitialized returns true if Sentry has been configured and a client exists.
func IsInitialized() bool {
	hub := sentry.CurrentHub()
	return hub != nil && hub.Client() != nil
}

// HubFromContext returns the Sentry hub from context, or the global hub if not found.
// Safe to call even if Sentry is not initialized.
func HubFromContext(ctx context.Context) *sentry.Hub {
	if sentry.HasHubOnContext(ctx) {
		return sentry.GetHubFromContext(ctx)
	}
	return sentry.CurrentHub()
}
