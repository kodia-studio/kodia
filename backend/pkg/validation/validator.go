package validation

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ContextRule is a validation function that requires context (e.g. database access).
type ContextRule func(ctx context.Context, fieldValue string) error

// Validator wraps go-playground/validator with Kodia-specific configuration.
type Validator struct {
	v        *validator.Validate
	ctxRules map[string]ContextRule
}

// New creates a Validator with all custom rules and JSON field name support.
func New() *Validator {
	v := validator.New()
	// Use JSON tag name instead of struct field name in error messages.
	// This makes errors show "email" instead of "Email".
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	registerCustomRules(v)
	return &Validator{v: v}
}

// Engine returns the underlying *validator.Validate for direct use if needed.
func (vl *Validator) Engine() *validator.Validate {
	return vl.v
}

// Struct validates a struct using its validate tags (sync validation only).
func (vl *Validator) Struct(s any) error {
	return vl.v.Struct(s)
}

// RegisterContextRule registers a rule that needs context (e.g. DB uniqueness check).
// Rules registered here are run by ValidateWithContext after sync validation passes.
// In DTO validate tags, prefix the rule name with "ctx:" (e.g. validate:"required,email,ctx:unique_email").
func (vl *Validator) RegisterContextRule(name string, fn ContextRule) {
	if vl.ctxRules == nil {
		vl.ctxRules = make(map[string]ContextRule)
	}
	vl.ctxRules[name] = fn
}

// ValidateWithContext runs sync validation then context rules.
// Returns either a validator.ValidationErrors (sync) or ContextValidationErrors (async).
func (vl *Validator) ValidateWithContext(ctx context.Context, s any) error {
	// 1. Sync validation (go-playground)
	if err := vl.v.Struct(s); err != nil {
		return err
	}
	// 2. Context rules
	return vl.runContextRules(ctx, s)
}

// runContextRules uses reflection to find fields with ctx:ruleName in validate tag.
func (vl *Validator) runContextRules(ctx context.Context, s any) error {
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	var ctxErrors []*ContextValidationError
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)
		if fv.Kind() != reflect.String {
			continue
		}

		tag := ft.Tag.Get("validate")
		fieldName := ft.Tag.Get("json")
		if fieldName == "" {
			fieldName = ft.Name
		}
		fieldName = strings.SplitN(fieldName, ",", 2)[0]

		for _, rule := range strings.Split(tag, ",") {
			rule = strings.TrimSpace(rule)
			if !strings.HasPrefix(rule, "ctx:") {
				continue
			}
			ruleName := rule[4:]
			fn, ok := vl.ctxRules[ruleName]
			if !ok {
				continue
			}
			if err := fn(ctx, fv.String()); err != nil {
				ctxErrors = append(ctxErrors, &ContextValidationError{
					Field:   fieldName,
					Message: err.Error(),
				})
			}
		}
	}
	if len(ctxErrors) > 0 {
		return &ContextValidationErrors{Errors: ctxErrors}
	}
	return nil
}

// ContextValidationError represents a single failed context validation rule.
type ContextValidationError struct {
	Field   string
	Message string
}

// ContextValidationErrors is the error type returned by ValidateWithContext.
type ContextValidationErrors struct {
	Errors []*ContextValidationError
}

func (e *ContextValidationErrors) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, ce := range e.Errors {
		msgs[i] = ce.Field + ": " + ce.Message
	}
	return strings.Join(msgs, "; ")
}

// FormatContextErrors converts ContextValidationErrors to map[field][]string.
func FormatContextErrors(err error) map[string][]string {
	var ce *ContextValidationErrors
	if !errors.As(err, &ce) {
		return map[string][]string{"error": {err.Error()}}
	}
	out := make(map[string][]string)
	for _, e := range ce.Errors {
		out[e.Field] = append(out[e.Field], e.Message)
	}
	return out
}

// FormatErrors converts go-playground validation errors to a user-friendly map[field][]string.
func FormatErrors(err error) map[string][]string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return map[string][]string{"error": {err.Error()}}
	}

	errs := make(map[string][]string)
	for _, fe := range ve {
		field := fe.Field()
		var msg string
		switch fe.Tag() {
		case "required":
			msg = field + " is required"
		case "email":
			msg = field + " must be a valid email address"
		case "min":
			msg = field + " is too short (min " + fe.Param() + " chars)"
		case "max":
			msg = field + " is too long (max " + fe.Param() + " chars)"
		case "len":
			msg = field + " must be exactly " + fe.Param() + " characters"
		case "url":
			msg = field + " must be a valid URL"
		case "uuid4":
			msg = field + " must be a valid UUID"
		case "strong_password":
			msg = field + " must contain uppercase, lowercase, number, and symbol"
		case "phone":
			msg = field + " must be a valid phone number"
		case "alpha_space":
			msg = field + " must contain only letters and spaces"
		case "no_html":
			msg = field + " must not contain HTML tags"
		default:
			msg = field + " is invalid (" + fe.Tag() + ")"
		}
		errs[field] = append(errs[field], msg)
	}
	return errs
}
