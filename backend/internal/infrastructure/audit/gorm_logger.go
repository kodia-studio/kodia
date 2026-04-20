package audit

import (
	"time"

	"github.com/kodia-studio/kodia/pkg/audit"
	"gorm.io/gorm"
)

type GormLogger struct {
	db *gorm.DB
}

// NewGormLogger creates a new audit logger that persists to database via GORM.
func NewGormLogger(db *gorm.DB) *GormLogger {
	// AutoMigrate the audit entry table
	_ = db.AutoMigrate(&audit.Entry{})
	
	return &GormLogger{db: db}
}

func (l *GormLogger) Log(entry audit.Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	return l.db.Create(&entry).Error
}

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
