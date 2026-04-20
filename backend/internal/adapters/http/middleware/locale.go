package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const ContextLocaleKey = "kodia:locale"

/**
 * LocaleMiddleware detects the user's preferred language.
 */
func LocaleMiddleware(defaultLocale string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check Header X-Locale
		locale := c.GetHeader("X-Locale")

		// 2. Check Query param 'lang'
		if locale == "" {
			locale = c.Query("lang")
		}

		// 3. Check Accept-Language header
		if locale == "" {
			acceptLang := c.GetHeader("Accept-Language")
			if acceptLang != "" {
				tags, _, err := language.ParseAcceptLanguage(acceptLang)
				if err == nil && len(tags) > 0 {
					locale = tags[0].String()
				}
			}
		}

		// Fallback to default
		if locale == "" {
			locale = defaultLocale
		}

		// Store in context
		c.Set(ContextLocaleKey, locale)
		
		c.Next()
	}
}

/**
 * GetLocale helper to retrieve the locale from context.
 */
func GetLocale(c *gin.Context) string {
	if locale, exists := c.Get(ContextLocaleKey); exists {
		return locale.(string)
	}
	return "en"
}
