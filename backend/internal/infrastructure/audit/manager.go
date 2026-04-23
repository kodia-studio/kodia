package audit

import (
	"github.com/kodia-studio/kodia/pkg/audit"
	"go.uber.org/zap"
)

// Manager coordinates multiple audit loggers (sinks).
type Manager struct {
	loggers []audit.Logger
	log     *zap.Logger // Internal logger for errors and structured fallback
}

// NewManager creates a new Audit Manager.
func NewManager(log *zap.Logger) *Manager {
	return &Manager{
		loggers: make([]audit.Logger, 0),
		log:     log,
	}
}

// AddLogger adds a new audit logger sink.
func (m *Manager) AddLogger(l audit.Logger) {
	m.loggers = append(m.loggers, l)
}

// Log broadcasts the audit entry to all registered sinks.
func (m *Manager) Log(entry audit.Entry) error {
	// Always log to structured log (Zap) as an external sink/fallback
	m.log.Info("Audit Log",
		zap.String("actor_id", entry.ActorID),
		zap.String("actor", entry.Actor),
		zap.String("action", string(entry.Action)),
		zap.String("resource", entry.Resource),
		zap.String("status", string(entry.Status)),
		zap.String("ip", entry.IP),
	)

	for _, l := range m.loggers {
		if err := l.Log(entry); err != nil {
			m.log.Error("Failed to record audit log to sink", zap.Error(err))
		}
	}
	return nil
}

// LogAction broadcasts the action to all registered sinks.
func (m *Manager) LogAction(actorID, actor, resource string, action audit.Action, status audit.Status, details string, ip, ua string) error {
	entry := audit.Entry{
		ActorID:   actorID,
		Actor:     actor,
		Action:    action,
		Resource:  resource,
		Details:   details,
		Status:    status,
		IP:        ip,
		UserAgent: ua,
	}
	return m.Log(entry)
}
