package models

import (
	"time"
)

// DocumentTemplate represents a template for documents, including metadata and relationships.
type DocumentTemplate struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	TitleEn        string     `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr        string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	DocumentTypeID uint       `json:"document_type_id" binding:"required" gorm:"not null"`
	Status         string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(45);default:'active';not null;check:status IN ('active', 'inactive')"`
	CreatedByID    uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedByID    *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	DocumentType  *DocumentType             `gorm:"foreignKey:DocumentTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *DocumentTemplateDocument `gorm:"foreignKey:DocumentTemplateID"`
	CreatedByUser *User                     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User                     `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentTemplate) TableName() string {
	return "document_templates"
}
