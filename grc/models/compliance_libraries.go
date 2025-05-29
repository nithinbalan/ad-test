package models

import "time"

type ComplianceLibrary struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Icon                   string    `json:"icon" gorm:"type:nvarchar(255)"`
	TitleEn                string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn          string    `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr          string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	IssuedBy               string    `json:"issued_by" gorm:"type:nvarchar(255)"`
	Classification         string    `json:"classification" gorm:"type:nvarchar(255)"`
	DocumentClassification string    `json:"document_classification" gorm:"type:nvarchar(255)"`
	Disclaimer             string    `json:"disclaimer" gorm:"type:nvarchar(max)"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime"`

	Versions []ComplianceLibraryVersion `json:"versions" gorm:"foreignKey:ComplianceLibraryID;constraint:OnDelete:CASCADE"`
}

func (ComplianceLibrary) TableName() string {
	return "compliance_libraries"
}
