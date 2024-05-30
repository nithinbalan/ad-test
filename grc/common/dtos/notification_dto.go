package dtos

import "jethur/grc/models"

type NotificationDTO struct {
	UserID           uint        `json:"user_id"`
	ModuleIdentifier string      `json:"module_id"`
	Icon             models.Icon `json:"icon"`
	MessageEn        string      `json:"message_en"`
	MessageAr        string      `json:"message_ar"`
	RedirectURL      string      `json:"redirect_url"`
	CreatedBy        uint        `json:"created_by"`
}
