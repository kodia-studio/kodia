package validation

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func registerCustomRules(v *validator.Validate) {
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("alpha_space", validateAlphaSpace)
	v.RegisterValidation("no_html", validateNoHTML)
}

// validateStrongPassword checks for uppercase, lowercase, digit, and symbol.
func validateStrongPassword(fl validator.FieldLevel) bool {
	p := fl.Field().String()
	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, c := range p {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSymbol = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSymbol
}

// validatePhone checks for E.164 phone format (flexible).
var phoneRegex = regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// validateAlphaSpace allows only letters and spaces.
func validateAlphaSpace(fl validator.FieldLevel) bool {
	for _, c := range fl.Field().String() {
		if !unicode.IsLetter(c) && c != ' ' {
			return false
		}
	}
	return true
}

// validateNoHTML rejects strings with < or > characters.
func validateNoHTML(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return !strings.Contains(s, "<") && !strings.Contains(s, ">")
}
