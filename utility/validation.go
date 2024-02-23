package utility

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExtractValidationErrors(c *gin.Context, err error) []string {
	var errors []string

	localizer, exists := c.Get("localizer")
	if !exists {
		return []string{"An error occurred while processing your request."}
	}

	loc, ok := localizer.(*i18n.Localizer)
	if !ok {
		return []string{"An error occurred while processing your request."}
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			var errorMessage string
			fieldName := fe.Field()

			switch fe.Tag() {
			case "required":
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_required",
					TemplateData: map[string]interface{}{
						"Field": fieldName,
					},
				})
			case "email":
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_email",
					TemplateData: map[string]interface{}{
						"Field": fieldName,
					},
				})
			default:
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_invalid",
					TemplateData: map[string]interface{}{
						"Field": fieldName,
					},
				})
			}

			errors = append(errors, errorMessage)
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}
