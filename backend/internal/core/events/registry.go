package events

import (
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// Registry maps event names to their listeners.
// This is the central source of truth for events in Kodia.
var Registry = map[string][]ports.Listener{
	// --- Listener Registration Start ---
	// "UserRegistered": {&listeners.SendWelcomeEmail{}, &listeners.LogAudit{}},
	// --- Listener Registration End ---
}

// RegisterEvents initializes the dispatcher with all listeners from the registry.
func RegisterEvents(dispatcher ports.EventDispatcher) {
	for eventName, listeners := range Registry {
		dispatcher.Register(eventName, listeners...)
	}
}
