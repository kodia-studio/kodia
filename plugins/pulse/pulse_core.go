package pulse

import (
	"go.uber.org/zap/zapcore"
)

type Core struct {
	zapcore.LevelEnabler
	enc     zapcore.Encoder
	manager *Manager
}

func NewCore(manager *Manager, level zapcore.Level) *Core {
	return &Core{
		LevelEnabler: level,
		enc:          zapcore.NewJSONEncoder(zapcore.NewProductionEncoderConfig()),
		manager:      manager,
	}
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	return c
}

func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Level >= zapcore.WarnLevel {
		c.manager.Log(entry.Level.String(), entry.Message)
	}
	return nil
}

func (c *Core) Sync() error {
	return nil
}
