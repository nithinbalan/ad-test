package utility

import (
	"encoding/json"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func LoadBundle(rootPath string) *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load message files for each language
	arabicLocalePath := filepath.Join(rootPath, "grc", "common", "locales", "ar", "error.ar.json")
	englishLocalePath := filepath.Join(rootPath, "grc", "common", "locales", "en", "error.en.json")

	bundle.MustLoadMessageFile(arabicLocalePath)
	bundle.MustLoadMessageFile(englishLocalePath)

	return bundle
}

// This function attempts to localize a message ID based on the context's localizer.
// If the localizer is not found or the localization fails, it returns a default message.
func LocalizeMessage(c *gin.Context, messageID string, defaultMessage string) string {
	// Attempt to get the localizer from the context
	localizerInterface, exists := c.Get("localizer")
	if !exists {
		return defaultMessage
	}

	// Assert the localizer to the correct type
	localizer, ok := localizerInterface.(*i18n.Localizer)
	if !ok {
		return defaultMessage
	}

	// Localize the error message
	localizedMessage, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return defaultMessage
	}

	return localizedMessage
}
