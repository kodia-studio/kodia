package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with Kodia-specific configuration.
type Validator struct {
	v *validator.Validate
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

// Struct validates a struct using its validate tags.
func (vl *Validator) Struct(s any) error {
	return vl.v.Struct(s)
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
