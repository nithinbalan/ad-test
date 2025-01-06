package dtos

import "time"

type EmailTemplateSettingsFormDTO struct {
	HeaderImageURL string `json:"header_image_url"`
	SupportText    string `json:"support_text" binding:"required"`
	ClosingText    string `json:"closing_text" binding:"required"`
	FooterLogoURL  string `json:"footer_logo_url"`
	FooterAddress  string `json:"footer_address" binding:"required"`
	FooterEmail    string `json:"footer_email" binding:"required"`
}

type EmailTemplateSettingsDTO struct {
	ID             int        `json:"id"`
	HeaderImageURL string     `json:"header_image_url"`
	SupportText    string     `json:"support_text"`
	ClosingText    string     `json:"closing_text"`
	FooterLogoURL  string     `json:"footer_logo_url"`
	FooterAddress  string     `json:"footer_address"`
	FooterEmail    string     `json:"footer_email"`
	CreatedByID    uint       `json:"created_by_id"`
	UpdatedByID    *uint      `json:"updated_by_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
