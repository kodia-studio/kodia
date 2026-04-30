package binder

import (
	"errors"

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

// Bind automatically binds the request JSON to the destination struct,
// sanitizes it, and performs validation. If it fails, it sends a standardized response
// and returns an error to stop execution in the handler.
// Sends 400 for JSON binding errors, 422 for validation errors.
func (b *Binding) Bind(c *gin.Context, dest interface{}) error {
	// 1. Bind JSON
	if err := c.ShouldBindJSON(dest); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return err
	}

	// 2. Sanitize
	_ = validation.Sanitize(dest)

	// 3. Validate Struct
	if err := b.validator.Struct(dest); err != nil {
		res := validation.FormatErrors(err)
		response.UnprocessableEntity(c, "Validation failed", res)
		return err
	}

	return nil
}

// BindCtx automatically binds the request JSON to the destination struct,
// sanitizes it, and performs both sync and async (context-aware) validation.
// Sends 400 for JSON binding errors, 422 for validation errors.
func (b *Binding) BindCtx(c *gin.Context, dest interface{}) error {
	// 1. Bind JSON
	if err := c.ShouldBindJSON(dest); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return err
	}

	// 2. Sanitize
	_ = validation.Sanitize(dest)

	// 3. Validate with context
	if err := b.validator.ValidateWithContext(c.Request.Context(), dest); err != nil {
		var errs map[string][]string
		var ce *validation.ContextValidationErrors
		if errors.As(err, &ce) {
			errs = validation.FormatContextErrors(err)
		} else {
			errs = validation.FormatErrors(err)
		}
		response.UnprocessableEntity(c, "Validation failed", errs)
		return err
	}

	return nil
}

// Global instance helper
var binder = New()

// Bind is a global shortcut for the Binding.Bind method.
func Bind(c *gin.Context, dest interface{}) error {
	return binder.Bind(c, dest)
}

// BindCtx is a global shortcut for the Binding.BindCtx method.
func BindCtx(c *gin.Context, dest interface{}) error {
	return binder.BindCtx(c, dest)
}
