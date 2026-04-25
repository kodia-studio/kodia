// Package events provides a public-facing event system for Kodia Framework.
// This package wraps the internal event dispatcher and provides developer-friendly types.
package events

import (
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// BaseEvent is the convenience base type for all user-defined domain events.
// Implement this to create custom events, or use it directly with Emit().
type BaseEvent struct {
	name    string
	payload interface{}
}

// NewEvent creates a new event with the given name and payload.
func NewEvent(name string, payload interface{}) *BaseEvent {
	return &BaseEvent{
		name:    name,
		payload: payload,
	}
}

// Name returns the event name.
func (e *BaseEvent) Name() string {
	return e.name
}

// Payload returns the event payload.
func (e *BaseEvent) Payload() interface{} {
	return e.payload
}

// Ensure BaseEvent implements ports.Event
var _ ports.Event = (*BaseEvent)(nil)
