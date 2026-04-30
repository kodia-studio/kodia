package validation

import (
	"html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Sanitize walks the struct fields of v (must be a pointer to struct)
// and applies sanitize: tag filters to every string or *string field.
// Filters are applied left to right in the order listed in the tag.
//
// Available filters:
//
//	trim         - strings.TrimSpace
//	lowercase    - strings.ToLower
//	uppercase    - strings.ToUpper
//	escape_html  - html.EscapeString (& → &amp; etc.)
//	strip_html   - strip all <...> HTML tags
//	slug         - lowercase, replace spaces/special chars with "-"
//	no_spaces    - remove all whitespace
//	alphanumeric - keep only letters and digits
//	normalize    - Unicode NFC normalization
//	truncate:N   - hard-truncate to N runes
func Sanitize(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return nil // not a struct pointer — silently skip
	}
	return sanitizeStruct(rv.Elem())
}

func sanitizeStruct(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		// Recurse into nested structs
		if fv.Kind() == reflect.Struct {
			if err := sanitizeStruct(fv); err != nil {
				return err
			}
			continue
		}
		if fv.Kind() == reflect.Ptr && !fv.IsNil() && fv.Elem().Kind() == reflect.Struct {
			if err := sanitizeStruct(fv.Elem()); err != nil {
				return err
			}
			continue
		}

		tag := ft.Tag.Get("sanitize")
		if tag == "" {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			if fv.CanSet() {
				fv.SetString(applyFilters(fv.String(), tag))
			}
		case reflect.Ptr:
			if !fv.IsNil() && fv.Elem().Kind() == reflect.String && fv.Elem().CanSet() {
				fv.Elem().SetString(applyFilters(fv.Elem().String(), tag))
			}
		}
	}
	return nil
}

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func applyFilters(s, tag string) string {
	for _, filter := range strings.Split(tag, ",") {
		filter = strings.TrimSpace(filter)
		switch {
		case filter == "trim":
			s = strings.TrimSpace(s)
		case filter == "lowercase":
			s = strings.ToLower(s)
		case filter == "uppercase":
			s = strings.ToUpper(s)
		case filter == "escape_html":
			s = html.EscapeString(s)
		case filter == "strip_html":
			s = htmlTagRe.ReplaceAllString(s, "")
		case filter == "slug":
			s = strings.ToLower(s)
			s = strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					return r
				}
				return '-'
			}, s)
			s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
			s = strings.Trim(s, "-")
		case filter == "no_spaces":
			s = strings.Join(strings.Fields(s), "")
		case filter == "alphanumeric":
			s = strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					return r
				}
				return -1
			}, s)
		case filter == "normalize":
			s = norm.NFC.String(s)
		case strings.HasPrefix(filter, "truncate:"):
			n, err := strconv.Atoi(filter[9:])
			if err == nil {
				runes := []rune(s)
				if len(runes) > n {
					s = string(runes[:n])
				}
			}
		}
	}
	return s
}
