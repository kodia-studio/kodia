package kodia

import (
	"sync"
)

/**
 * Kodia Hook System 🐨🔌
 * A high-performance, thread-safe event dispatcher for the plugin ecosystem.
 */

// HookCallback is the signature for functions listening to hooks.
type HookCallback func(data any)

// HookManager manages registration and dispatching of event hooks.
type HookManager struct {
	mu    sync.RWMutex
	hooks map[string][]HookCallback
}

// NewHookManager initializes a new hook manager.
func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[string][]HookCallback),
	}
}

// Listen registers a callback for a specific event.
func (m *HookManager) Listen(event string, callback HookCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks[event] = append(m.hooks[event], callback)
}

// Dispatch executes all callbacks registered for a specific event.
func (m *HookManager) Dispatch(event string, data any) {
	m.mu.RLock()
	callbacks, exists := m.hooks[event]
	m.mu.RUnlock()

	if !exists {
		return
	}

	for _, callback := range callbacks {
		// We execute callbacks synchronously for predictable flow,
		// but plugins can spawn goroutines if they need async behavior.
		callback(data)
	}
}

// HasListeners checks if an event has any registered callbacks.
func (m *HookManager) HasListeners(event string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.hooks[event]
	return exists
}
