package binder

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/kodia-studio/kodia/pkg/validation"
)

// Binding provides the entry point for request binding and validation.
type Binding struct {
	validator *validation.Validator
}

// New creates a new Binding instance.
func New() *Binding {
	return &Binding{
		validator: validation.New(),
	}
}

// Bind automatically binds the request JSON to the destination struct
// and performs validation. If it fails, it sends a standardized 400 response
// and returns an error to stop execution in the handler.
func (b *Binding) Bind(c *gin.Context, dest interface{}) error {
	// 1. Bind JSON
	if err := c.ShouldBindJSON(dest); err != nil {
		response.BadRequest(c, "Invalid request payload", map[string]string{"error": err.Error()})
		return err
	}

	// 2. Validate Struct
	if err := b.validator.Struct(dest); err != nil {
		res := validation.FormatErrors(err)
		response.BadRequest(c, "Validation failed", res)
		return err
	}

	return nil
}

// Global instance helper
var binder = New()

// Bind is a global shortcut for the Binder.Bind method.
func Bind(c *gin.Context, dest interface{}) error {
	return binder.Bind(c, dest)
}
