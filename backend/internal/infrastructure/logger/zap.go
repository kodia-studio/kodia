// Package logger provides a structured, production-ready logger for Kodia Framework.
// Uses uber-go/zap under the hood.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates and returns a new zap.Logger.
// In production, it uses JSON encoding with INFO level.
// In development, it uses console encoding with DEBUG level.
func New(isDevelopment bool) (*zap.Logger, error) {
	var cfg zap.Config

	if isDevelopment {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	return cfg.Build()
}
