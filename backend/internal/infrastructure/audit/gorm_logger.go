package audit

import (
	"time"

	"github.com/kodia-studio/kodia/pkg/audit"
	"gorm.io/gorm"
)

// GormLogger implements audit.Logger using a database via GORM.
type GormLogger struct {
	db *gorm.DB
}

// NewGormLogger creates a new GormLogger instance.
func NewGormLogger(db *gorm.DB) *GormLogger {
	return &GormLogger{db: db}
}

// Log records an audit entry to the database.
func (l *GormLogger) Log(entry audit.Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	return l.db.Create(&entry).Error
}

// LogAction is a convenience method to record an action with its details.
func (l *GormLogger) LogAction(actorID, actor, resource string, action audit.Action, status audit.Status, details string, ip, ua string) error {
	entry := audit.Entry{
		Timestamp: time.Now(),
		ActorID:   actorID,
		Actor:     actor,
		Action:    action,
		Resource:  resource,
		Details:   details,
		Status:    status,
		IP:        ip,
		UserAgent: ua,
	}
	return l.Log(entry)
}
