package domain

import "errors"

// Domain errors — use these throughout the application for consistent error handling.
// Never return raw database errors to the handler layer; map them to domain errors.

var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when trying to create a duplicate resource.
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidCredentials is returned when authentication fails.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUnauthorized is returned when the user is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission.
	ErrForbidden = errors.New("forbidden")

	// ErrTokenExpired is returned when a JWT or refresh token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenRevoked is returned when a refresh token has been revoked.
	ErrTokenRevoked = errors.New("token revoked")

	// ErrInactiveAccount is returned when a user tries to log in with a deactivated account.
	ErrInactiveAccount = errors.New("account is inactive")

	// ErrInvalidInput is returned for general invalid input that doesn't fall under validation.
	ErrInvalidInput = errors.New("invalid input")
)
