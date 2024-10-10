package models

import (
	"time"
)

type EmailTemplateSettings struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	HeaderImageURL string     `json:"header_image_url" binding:"required" gorm:"type:nvarchar(max);not null"`
	SupportText    string     `json:"support_text" binding:"required" gorm:"type:nvarchar(max);not null"`
	ClosingText    string     `json:"closing_text" binding:"required" gorm:"type:nvarchar(max);not null"`
	FooterLogoURL  string     `json:"footer_logo_url" binding:"required" gorm:"type:nvarchar(max);not null"`
	FooterAddress  string     `json:"footer_address" binding:"required" gorm:"type:nvarchar(max);not null"`
	FooterEmail    string     `json:"footer_email" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedByID    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID    *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime;not null"`
}

// TableName overrides the table name used by GORM
func (EmailTemplateSettings) TableName() string {
	return "email_template_settings"
}
