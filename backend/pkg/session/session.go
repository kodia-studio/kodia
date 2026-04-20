// Package session provides types for enterprise-grade session and device tracking.
package session

import (
	"time"
)

// Session represents an active user session/device.
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"index"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Device    string    `json:"device"` // Extracted from UA (e.g., iPhone, Chrome on Windows)
	Location  string    `json:"location"` // Optional: Country/City from IP
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen" gorm:"index"`
	IsExpired bool      `json:"is_expired" gorm:"index"`
}

// Store defines the interface for persisting and managing sessions.
type Store interface {
	Create(session *Session) error
	Get(id string) (*Session, error)
	GetByUserID(userID string) ([]Session, error)
	UpdateLastSeen(id string) error
	Revoke(id string) error
	RevokeAllForUser(userID string) error
}
