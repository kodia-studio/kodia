package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// LocaleMessages maps a validation tag to its human-readable message template.
// Use {field} and {param} as placeholders.
type LocaleMessages map[string]string

var messages = map[string]LocaleMessages{
	"en": {
		"required":        "{field} is required",
		"email":           "{field} must be a valid email address",
		"min":             "{field} is too short (min {param} chars)",
		"max":             "{field} is too long (max {param} chars)",
		"len":             "{field} must be exactly {param} characters",
		"url":             "{field} must be a valid URL",
		"uuid4":           "{field} must be a valid UUID",
		"strong_password": "{field} must contain uppercase, lowercase, number, and symbol",
		"phone":           "{field} must be a valid phone number",
		"alpha_space":     "{field} must contain only letters and spaces",
		"no_html":         "{field} must not contain HTML tags",
		"alphanumeric":    "{field} must contain only letters and digits",
		"slug":            "{field} must be a valid slug (lowercase, letters, digits, hyphens)",
		"date_format":     "{field} must match date format {param}",
	},
	"id": {
		"required":        "{field} wajib diisi",
		"email":           "{field} harus berupa alamat email yang valid",
		"min":             "{field} terlalu pendek (minimal {param} karakter)",
		"max":             "{field} terlalu panjang (maksimal {param} karakter)",
		"len":             "{field} harus tepat {param} karakter",
		"url":             "{field} harus berupa URL yang valid",
		"uuid4":           "{field} harus berupa UUID yang valid",
		"strong_password": "{field} harus mengandung huruf besar, kecil, angka, dan simbol",
		"phone":           "{field} harus berupa nomor telepon yang valid",
		"alpha_space":     "{field} hanya boleh mengandung huruf dan spasi",
		"no_html":         "{field} tidak boleh mengandung tag HTML",
		"alphanumeric":    "{field} hanya boleh mengandung huruf dan angka",
		"slug":            "{field} harus berupa slug yang valid",
		"date_format":     "{field} harus sesuai format tanggal {param}",
	},
}

// FormatErrorsLocale formats validation errors using the given locale (e.g. "en", "id").
// Falls back to English for unknown locales or unknown tags.
func FormatErrorsLocale(err error, locale string) map[string][]string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return map[string][]string{"error": {err.Error()}}
	}

	msgs, ok := messages[locale]
	if !ok {
		msgs = messages["en"]
	}

	errs := make(map[string][]string)
	for _, fe := range ve {
		field := fe.Field()
		tmpl, ok := msgs[fe.Tag()]
		if !ok {
			tmpl = messages["en"][fe.Tag()]
		}
		if tmpl == "" {
			tmpl = "{field} is invalid (" + fe.Tag() + ")"
		}
		msg := strings.NewReplacer("{field}", field, "{param}", fe.Param()).Replace(tmpl)
		errs[field] = append(errs[field], msg)
	}
	return errs
}

// RegisterLocale adds or overrides messages for a locale.
// Use to add custom rules or entirely new languages.
func RegisterLocale(locale string, msgs LocaleMessages) {
	if existing, ok := messages[locale]; ok {
		for k, v := range msgs {
			existing[k] = v
		}
	} else {
		messages[locale] = msgs
	}
}
