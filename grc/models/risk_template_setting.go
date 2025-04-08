package models

import (
	"time"
)

type RiskTemplateSetting struct {
	ID                   uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	LogoImagePath        *string    `json:"logo_image_path" gorm:"type:varchar(500);default:null"`
	CoverPageBGImagePath *string    `json:"cover_page_bg_image_path" gorm:"type:varchar(500);default:null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskTemplateSetting to `audit_template_settings`
func (RiskTemplateSetting) TableName() string {
	return "risk_template_settings"
}
