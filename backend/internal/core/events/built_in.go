package events

import (
	"time"
)

// --- User Domain Events ---

// UserRegisteredEvent is dispatched when a new user registers.
type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserRegisteredEvent) Name() string {
	return "user.registered"
}

func (e UserRegisteredEvent) Payload() interface{} {
	return e
}

// UserLoggedInEvent is dispatched when a user successfully logs in.
type UserLoggedInEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserLoggedInEvent) Name() string {
	return "user.logged_in"
}

func (e UserLoggedInEvent) Payload() interface{} {
	return e
}

// PasswordChangedEvent is dispatched when a user changes their password.
type PasswordChangedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PasswordChangedEvent) Name() string {
	return "password.changed"
}

func (e PasswordChangedEvent) Payload() interface{} {
	return e
}

// PasswordResetEvent is dispatched when a password reset is requested or completed.
type PasswordResetEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // "requested", "completed"
	Timestamp time.Time `json:"timestamp"`
}

func (e PasswordResetEvent) Name() string {
	return "password.reset"
}

func (e PasswordResetEvent) Payload() interface{} {
	return e
}

// TwoFactorEnabledEvent is dispatched when 2FA is enabled by a user.
type TwoFactorEnabledEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e TwoFactorEnabledEvent) Name() string {
	return "two_factor.enabled"
}

func (e TwoFactorEnabledEvent) Payload() interface{} {
	return e
}

// TwoFactorDisabledEvent is dispatched when 2FA is disabled by a user.
type TwoFactorDisabledEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e TwoFactorDisabledEvent) Name() string {
	return "two_factor.disabled"
}

func (e TwoFactorDisabledEvent) Payload() interface{} {
	return e
}

// --- Role/Permission Events ---

// RoleAssignedEvent is dispatched when a role is assigned to a user.
type RoleAssignedEvent struct {
	UserID    string    `json:"user_id"`
	RoleName  string    `json:"role_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e RoleAssignedEvent) Name() string {
	return "role.assigned"
}

func (e RoleAssignedEvent) Payload() interface{} {
	return e
}

// RoleRevokedEvent is dispatched when a role is revoked from a user.
type RoleRevokedEvent struct {
	UserID    string    `json:"user_id"`
	RoleName  string    `json:"role_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e RoleRevokedEvent) Name() string {
	return "role.revoked"
}

func (e RoleRevokedEvent) Payload() interface{} {
	return e
}

// --- Generic Email Verification Events ---

// EmailVerificationRequestedEvent is dispatched when email verification is requested.
type EmailVerificationRequestedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e EmailVerificationRequestedEvent) Name() string {
	return "email.verification_requested"
}

func (e EmailVerificationRequestedEvent) Payload() interface{} {
	return e
}

// EmailVerifiedEvent is dispatched when an email is verified.
type EmailVerifiedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e EmailVerifiedEvent) Name() string {
	return "email.verified"
}

func (e EmailVerifiedEvent) Payload() interface{} {
	return e
}
