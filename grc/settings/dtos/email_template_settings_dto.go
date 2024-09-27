package dtos

type EmailTemplateSettingsFormDTO struct {
	HeaderImageURL string `json:"header_image_url" binding:"required"`
	SupportText    string `json:"support_text" binding:"required"`
	ClosingText    string `json:"closing_text" binding:"required"`
	FooterLogoURL  string `json:"footer_logo_url" binding:"required"`
	FooterAddress  string `json:"footer_address" binding:"required"`
	FooterEmail    string `json:"footer_email" binding:"required"`
}

// Use this DTO for creating new settings or updating existing ones, depending on your application's needs.
