package validation

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/response"
)

const validatedKey = "kodia_validated_request"

// BindAndValidate binds JSON body to req and validates it.
// Sends 422 Unprocessable Entity with formatted errors if validation fails.
// Applies sanitization before validation.
// Returns true if binding and validation succeed, false otherwise.
// If false is returned, the error response has already been sent.
func BindAndValidate(c *gin.Context, vl *Validator, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return false
	}
	_ = Sanitize(req)
	if err := vl.Struct(req); err != nil {
		response.UnprocessableEntity(c, "Validation failed", FormatErrors(err))
		return false
	}
	return true
}

// BindAndValidateCtx binds JSON body to req and validates it with context-aware rules.
// Sends 422 Unprocessable Entity with formatted errors if validation fails.
// Applies sanitization before validation.
// Returns true if binding and validation succeed, false otherwise.
// If false is returned, the error response has already been sent.
func BindAndValidateCtx(c *gin.Context, vl *Validator, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return false
	}
	_ = Sanitize(req)
	if err := vl.ValidateWithContext(c.Request.Context(), req); err != nil {
		var errs map[string][]string
		var ce *ContextValidationErrors
		if errors.As(err, &ce) {
			errs = FormatContextErrors(err)
		} else {
			errs = FormatErrors(err)
		}
		response.UnprocessableEntity(c, "Validation failed", errs)
		return false
	}
	return true
}

// Middleware returns a gin.HandlerFunc that binds and validates T before calling the handler.
// If validation fails, responds with 422 and aborts. Otherwise, injects the validated request into context.
func Middleware[T any](vl *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if !BindAndValidate(c, vl, &req) {
			c.Abort()
			return
		}
		c.Set(validatedKey, req)
		c.Next()
	}
}

// Get retrieves the validated request injected by Middleware[T].
func Get[T any](c *gin.Context) T {
	return c.MustGet(validatedKey).(T)
}
