package ports

import (
	"context"
)

// Event represents a system event.
type Event interface {
	Name() string
	Payload() interface{}
}

// Listener represents a handler for an event.
type Listener interface {
	Handle(ctx context.Context, event Event) error
}

// AsyncListener marks a listener that should be queued.
type AsyncListener interface {
	Listener
	ShouldQueue() bool
}

// EventDispatcher manages event dispatching.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event) error
	Register(eventName string, listeners ...Listener)
}
