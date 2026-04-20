package observability

import (
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

// PulseCore is a zapcore.Core that streams logs to PulseManager.
type PulseCore struct {
	zapcore.LevelEnabler
	enc     zapcore.Encoder
	manager *PulseManager
}

// NewPulseCore creates a new PulseCore.
func NewPulseCore(manager *PulseManager, level zapcore.Level) *PulseCore {
	return &PulseCore{
		LevelEnabler: level,
		enc: zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		manager:      manager,
	}
}

// With returns a new core with the given fields.
func (c *PulseCore) With(fields []zapcore.Field) zapcore.Core {
	return c
}

// Check determines whether the entry should be logged.
func (c *PulseCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

// Write streams the log entry to PulseManager.
func (c *PulseCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Only stream Warning and Error by default as per user request
	if entry.Level >= zapcore.WarnLevel {
		c.manager.Log(entry.Level.String(), entry.Message)
	}
	return nil
}

// Sync is a no-op.
func (c *PulseCore) Sync() error {
	return nil
}
