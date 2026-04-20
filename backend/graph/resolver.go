package graph

import (
	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthService ports.AuthService
	UserService ports.UserService
	Log         *zap.Logger
}
