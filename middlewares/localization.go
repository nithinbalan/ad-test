package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func LocalizerMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		if accept == "" {
			accept = "en"
		}
		tags, _, err := language.ParseAcceptLanguage(accept)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Accept-Language header"})
			c.Abort()
			return
		}

		// Choose the most appropriate language from the Accept-Language header
		// and the languages loaded in the bundle.
		lang := bundle.LanguageTags()
		matcher := language.NewMatcher(lang)
		tag, _, _ := matcher.Match(tags...)

		localizer := i18n.NewLocalizer(bundle, tag.String())

		// Store the localizer in the context
		c.Set("localizer", localizer)
		c.Next()
	}
}
