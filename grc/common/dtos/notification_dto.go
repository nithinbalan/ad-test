package dtos

import "jethur/grc/models"

type DBNotificationDTO struct {
	UserID           uint        `json:"user_id"`
	ModuleIdentifier string      `json:"module_id"`
	Icon             models.Icon `json:"icon"`
	MessageEn        string      `json:"message_en"`
	MessageAr        string      `json:"message_ar"`
	RedirectURL      string      `json:"redirect_url"`
	CreatedBy        uint        `json:"created_by"`
}

type EmailNotificationDTO struct {
	EmailTo              []string    `json:"email_to"`
	Subject              string      `json:"subject"`
	TemplateRelativePath string      `json:"template_relative_path"`
	EmailData            interface{} `json:"email_data"`
}
