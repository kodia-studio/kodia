// Package authsaml provides types and interfaces for SAML (Security Assertion Markup Language) support.
package authsaml

// Provider defines the interface for SAML identity provider integration.
type Provider interface {
	GetSSOURL() string
	ValidateAssertion(responseXML string) (*Identity, error)
	GenerateMetadata() (string, error)
}

// Identity represents the authenticated subject from a SAML assertion.
type Identity struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	Attributes map[string]string `json:"attributes"`
}

// Config holds the SAML service provider configuration.
type Config struct {
	EntityID    string
	SSOURL      string
	Certificate string
}
