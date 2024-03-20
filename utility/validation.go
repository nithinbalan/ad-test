package utility

import (
	"encoding/json"

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

	// Handle JSON syntax error
	if _, ok := err.(*json.SyntaxError); ok {
		errorMessage, _ := loc.Localize(&i18n.LocalizeConfig{MessageID: "invalid_json_format"})
		errors = append(errors, errorMessage)
		return errors
	}

	// Handle JSON unmarshal type error (when JSON type doesn't match struct field type)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		errorMessage, _ := loc.Localize(&i18n.LocalizeConfig{MessageID: "incorrect_json_field_type"})
		errors = append(errors, errorMessage)
		return errors
	}

	//Error from custom validation
	if ve, ok := err.(*ValidationError); ok {

		translatedFieldName, err := loc.Localize(&i18n.LocalizeConfig{MessageID: RemoveLangCodeSuffix(ve.field)})
		if err != nil {
			translatedFieldName = ve.field
		}

		switch ve.errType {
		case "unique":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "resource_already_exists",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "required":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "field_required",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		}

		return errors
	}

	//Validation errors from GIN
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			var errorMessage string
			translatedFieldName, err := loc.Localize(&i18n.LocalizeConfig{MessageID: RemoveLangCodeSuffix(fe.Field())})
			if err != nil {
				translatedFieldName = fe.Field()
			}
			switch fe.Tag() {
			case "required":
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_required",
					TemplateData: map[string]interface{}{
						"Field": translatedFieldName,
					},
				})
			case "email":
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_email",
					TemplateData: map[string]interface{}{
						"Field": translatedFieldName,
					},
				})
			default:
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_invalid",
					TemplateData: map[string]interface{}{
						"Field": translatedFieldName,
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
