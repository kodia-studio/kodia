package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

/**
 * Translator handles message lookups across multiple locales.
 */
type Translator struct {
	bundle *i18n.Bundle
}

func NewTranslator(defaultLang string) *Translator {
	bundle := i18n.NewBundle(language.Make(defaultLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return &Translator{bundle: bundle}
}

// LoadMessageFile loads a translation file from the given path.
func (t *Translator) LoadMessageFile(path string) error {
	_, err := t.bundle.LoadMessageFile(path)
	return err
}

// T translates a message with optional parameters.
func (t *Translator) T(locale string, messageID string, params map[string]interface{}) string {
	localizer := i18n.NewLocalizer(t.bundle, locale)
	
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: params,
	})

	if err != nil {
		return messageID // Fallback to ID if not found
	}

	return msg
}

// Bundle returns the underlying i18n bundle if needed.
func (t *Translator) Bundle() *i18n.Bundle {
	return t.bundle
}
