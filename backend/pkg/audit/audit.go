// Package audit provides types and interfaces for enterprise-grade audit logging in Kodia.
package audit

import (
	"time"
)

// Action represents the type of activity being logged.
type Action string

const (
	ActionCreate Action = "CREATE"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
	ActionLogin  Action = "LOGIN"
	ActionLogout Action = "LOGOUT"
	ActionExport Action = "EXPORT"
	ActionAuth2  Action = "AUTH_2FA"
)

// Status represents the outcome of the action.
type Status string

const (
	StatusSuccess Status = "SUCCESS"
	StatusFailure Status = "FAILURE"
)

// Entry represents a single audit log record.
type Entry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	ActorID   string    `json:"actor_id" gorm:"index"`
	Actor     string    `json:"actor"` // Email or Username
	Action    Action    `json:"action" gorm:"index"`
	Resource  string    `json:"resource" gorm:"index"`
	Details   string    `json:"details" gorm:"type:text"`
	Status    Status    `json:"status"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
}

// Logger defines the interface for recording audit entries.
type Logger interface {
	Log(entry Entry) error
	LogAction(actorID, actor, resource string, action Action, status Status, details string, ip, ua string) error
}
