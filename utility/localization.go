package utility

import (
	"encoding/json"
	"jethur/assets"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func LoadBundle(rootPath string) *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load message files for each language
	commonArErrorPath := "common/locales/ar/error.ar.json"
	commonArLabelsPath := "common/locales/ar/labels.ar.json"
	commonEnErrorPath := "common/locales/en/error.en.json"
	commonEnLabelsPath := "common/locales/en/labels.en.json"

	embeddedArError, err := assets.FS.ReadFile(commonArErrorPath)
	if err != nil {
		LogFatal(err.Error(), nil)
	}
	embeddedArLabels, err := assets.FS.ReadFile(commonArLabelsPath)
	if err != nil {
		LogFatal(err.Error(), nil)
	}
	embeddedEnError, err := assets.FS.ReadFile(commonEnErrorPath)
	if err != nil {
		LogFatal(err.Error(), nil)
	}
	embeddedEnLabels, err := assets.FS.ReadFile(commonEnLabelsPath)
	if err != nil {
		LogFatal(err.Error(), nil)
	}

	// Parse and load embedded JSON files into the bundle
	bundle.MustParseMessageFileBytes(embeddedArError, commonArErrorPath)
	bundle.MustParseMessageFileBytes(embeddedArLabels, commonArLabelsPath)
	bundle.MustParseMessageFileBytes(embeddedEnError, commonEnErrorPath)
	bundle.MustParseMessageFileBytes(embeddedEnLabels, commonEnLabelsPath)

	return bundle
}

func LocalizeMessage(c *gin.Context, messageID string, defaultMessage string) string {
	localizerInterface, exists := c.Get("localizer")
	if !exists {
		return defaultMessage
	}

	localizer, ok := localizerInterface.(*i18n.Localizer)
	if !ok {
		return defaultMessage
	}

	localizedMessage, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return defaultMessage
	}

	return localizedMessage
}
