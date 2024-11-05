package models

import "time"

type Integration struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IntegrationCategoryID uint       `json:"integration_category_id" binding:"required"`
	TitleEn               string     `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr               string     `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn         *string    `json:"description_en" gorm:"type:text"`
	DescriptionAr         *string    `json:"description_ar" gorm:"type:text"`
	LogoUrl               *string    `json:"logo_url" gorm:"type:text"`
	ApiBaseUrl            *string    `json:"api_base_url" gorm:"type:text"`
	IsEnable              bool       `json:"is_enable" gorm:"not null;type:tinyint;default:0" binding:"required"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	CreatedByUser *User                    `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                    `gorm:"foreignKey:UpdatedBYID"`
	Credentials   []*IntegrationCredential `gorm:"foreignKey:IntegrationID"`
}

// TableName overrides the table name used by Integration to `integrations`
func (Integration) TableName() string {
	return "integrations"
}
