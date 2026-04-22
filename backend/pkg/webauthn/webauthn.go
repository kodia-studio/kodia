package webauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

// Manager wraps the webauthn library with Kodia configuration.
type Manager struct {
	instance *webauthn.WebAuthn
}

// Config holds the WebAuthn specific configuration.
type Config struct {
	RPDisplayName string
	RPID          string
	RPOrigins     []string
}

// NewManager creates a new WebAuthn Manager.
func NewManager(cfg Config) (*Manager, error) {
	w, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.RPDisplayName,
		RPID:          cfg.RPID,
		RPOrigins:     cfg.RPOrigins,
	})
	if err != nil {
		return nil, err
	}

	return &Manager{instance: w}, nil
}

// Engine returns the underlying *webauthn.WebAuthn instance.
func (m *Manager) Engine() *webauthn.WebAuthn {
	return m.instance
}
