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
		case "not_exist":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "resource_does_not_exist",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_date":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_date",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "due_date_before_start_date":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "due_date_before_start_date",
			})
			errors = append(errors, errorMessage)
		case "invalid":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "resend_not_yet_allowed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "resend_not_yet_allowed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "otp_expired":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "otp_expired",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "session_expired":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "session_expired",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "password_complexity":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "password_complexity",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "password_mismatch":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "password_mismatch",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_otp":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_otp",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_credentials":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_credentials",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "decryption_failed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "decryption_failed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_file_type":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_file_type",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_word_file_type":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_word_file_type",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "duplicate":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "duplicate",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "empty":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "empty",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_clause_number":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_clause_number",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_image_type":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_image_type",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_profile_file_size":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_profile_file_size",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "already_associated":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "already_associated",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "matrix_exists":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "matrix_exists",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_header":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_header",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_data":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_data",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "edit_not_allowed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "edit_not_allowed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "distinct_review_user":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "distinct_review_user",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "user_exceed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "user_exceed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "module_not_allowed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "module_not_allowed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "default_criticality_matrix":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "default_criticality_matrix",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "same_organization_unit":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "same_organization_unit",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "ai_not_allowed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "ai_not_allowed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "invalid_range":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_range",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)
		case "out_of_range":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "out_of_range",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			if errorMessage == "" {
				errorMessage = ve.Error()
			}
			errors = append(errors, errorMessage)

		case "kri_maximum_susceptibility_percent_invalid_range":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "kri_maximum_susceptibility_percent_invalid_range",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			if errorMessage == "" {
				errorMessage = ve.Error()
			}
			errors = append(errors, errorMessage)

		case "kri_maximum_tolerable_percent_invalid_range":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "kri_maximum_tolerable_percent_invalid_range",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			if errorMessage == "" {
				errorMessage = ve.Error()
			}
			errors = append(errors, errorMessage)

		case "document_not_published":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "document_not_published",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "field_max_length":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "field_max_length",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "invalid_total_weightage":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_total_weightage",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "document_upload_disabled":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "document_upload_disabled",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "already_acknowledged":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "already_acknowledged",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "inherent_risk_analysis_incomplete":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "inherent_risk_analysis_incomplete",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "missing_participants":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "missing_participants",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "invalid_hierarchy":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "invalid_hierarchy",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "residual_risk_rating":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "residual_risk_rating",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "already_used":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "already_used",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "no_unit_head":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "no_unit_head",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "register_not_approved":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "register_not_approved",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "bucket_range":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "bucket_range",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "not_valid":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "not_valid",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "not_found":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "not_found",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "empty_data":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "empty_data",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "employee_analysise_required":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "employee_analysise_required",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_already_submitted":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_already_submitted",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_already_completed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_already_completed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_not_active":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_not_active",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "review_comment_workflow_not_active":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "review_comment_workflow_not_active",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_not_rejected":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_not_rejected",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_submission_required":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_submission_required",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "workflow_submission_requires_one":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_submission_requires_one",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "subsequent_stage_completed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "subsequent_stage_completed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "reporting_user_not_found":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "reporting_user_not_found",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "no_users_for_role":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "no_users_for_role",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "incident_closure_submission":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "incident_closure_submission",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "modify_requires_latest":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "modify_requires_latest",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "latest_must_be_failed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "latest_must_be_failed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "empty_result":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "empty_result",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "draft_exists":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "draft_exists",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "retest_already_completed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "retest_already_completed",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "already_published":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "already_published",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "retest_already_active":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "retest_already_active",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "no_requirements":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "no_requirements",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "version_published":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "version_published",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "test_end_before_start":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "test_end_before_start",
			})
			errors = append(errors, errorMessage)

		case "no_published_version":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "no_published_version",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "stage_change_not_allowed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "stage_change_not_allowed",
			})
			errors = append(errors, errorMessage)

		case "checklist_not_marked_as_exception":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "checklist_not_marked_as_exception",
			})
			errors = append(errors, errorMessage)

		case "workflow_already_approved":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "workflow_already_approved",
				TemplateData: map[string]interface{}{
					"Field": translatedFieldName,
				},
			})
			errors = append(errors, errorMessage)

		case "security_gate_checklist_not_answered":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "security_gate_checklist_not_answered",
			})
			errors = append(errors, errorMessage)

		case "security_gate_checklist_not_implemented":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "security_gate_checklist_not_implemented",
			})
			errors = append(errors, errorMessage)

		case "security_gate_exception_not_approved":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "security_gate_exception_not_approved",
			})
			errors = append(errors, errorMessage)

		case "security_gate_already_go_live":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "security_gate_already_go_live",
			})
			errors = append(errors, errorMessage)

		case "findings_not_closed":
			var errorMessage string
			errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID: "findings_not_closed",
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
			case "eqfield":
				errorMessage, _ = loc.Localize(&i18n.LocalizeConfig{
					MessageID: "field_confirm_password",
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
