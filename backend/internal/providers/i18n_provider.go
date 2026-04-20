package providers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kodia-studio/kodia/pkg/i18n"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
)

type I18nServiceProvider struct {
	translator *i18n.Translator
}

func NewI18nServiceProvider() *I18nServiceProvider {
	return &I18nServiceProvider{}
}

func (p *I18nServiceProvider) Name() string {
	return "kodia:i18n"
}

func (p *I18nServiceProvider) Register(app *kodia.App) error {
	defaultLang := app.Config.App.Locale
	if defaultLang == "" {
		defaultLang = "en"
	}

	translator := i18n.NewTranslator(defaultLang)
	p.translator = translator
	
	// Register in container
	app.Set("translator", translator)
	
	return nil
}

func (p *I18nServiceProvider) Boot(app *kodia.App) error {
	langDir := filepath.Join("resources", "lang")
	
	// Ensure directory exists
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		app.Log.Warn("i18n language directory not found, skipping message loading", zap.String("path", langDir))
		return nil
	}

	// Load all JSON files in resources/lang
	files, err := os.ReadDir(langDir)
	if err != nil {
		return fmt.Errorf("failed to read lang directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(langDir, file.Name())
			if err := p.translator.LoadMessageFile(filePath); err != nil {
				app.Log.Error("Failed to load translation file", zap.String("file", file.Name()), zap.Error(err))
			} else {
				app.Log.Debug("Loaded translation file", zap.String("file", file.Name()))
			}
		}
	}

	return nil
}
