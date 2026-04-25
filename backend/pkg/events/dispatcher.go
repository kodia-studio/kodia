package events

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/ports"
)

// Dispatcher is the public event bus for developers.
// It wraps the internal dispatcher with convenient aliases and methods.
type Dispatcher struct {
	internal ports.EventDispatcher
}

// NewDispatcher creates a new public Dispatcher wrapping the internal dispatcher.
func NewDispatcher(internal ports.EventDispatcher) *Dispatcher {
	return &Dispatcher{
		internal: internal,
	}
}

// Emit dispatches an event to all registered listeners.
// Listeners marked as async will be queued for background processing.
func (d *Dispatcher) Emit(ctx context.Context, event ports.Event) error {
	return d.internal.Dispatch(ctx, event)
}

// On registers one or more listeners for a specific event name.
func (d *Dispatcher) On(eventName string, listeners ...ports.Listener) {
	d.internal.Register(eventName, listeners...)
}

// EmitAsync is a helper method to emit an event and ensure async processing.
// Note: This method requires the event to be handled by AsyncListener implementations.
// For true fire-and-forget, listeners should implement ShouldQueue() to return true.
func (d *Dispatcher) EmitAsync(ctx context.Context, event ports.Event) error {
	return d.Emit(ctx, event)
}
