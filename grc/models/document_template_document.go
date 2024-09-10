package models

import (
	"time"
)

type DocumentTemplateDocument struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	DocumentTemplateID uint       `json:"document_template_id" gorm:"not null" binding:"required"`
	URL                string     `json:"url" gorm:"type:varchar(1000);not null" binding:"required"`
	Name               string     `json:"name" gorm:"type:nvarchar(255);not null" binding:"required"`
	Size               int        `json:"size" gorm:"not null" binding:"required"`
	Token              string     `json:"token" gorm:"type:varchar(255);not null" binding:"required"`
	OriginalName       string     `json:"original_name"`
	Type               string     `json:"type" gorm:"type:varchar(255);not null" binding:"required"`
	CreatedByID        uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedByID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	//Relation ships
	DocumentTemplate *DocumentTemplate `gorm:"foreignKey:DocumentTemplateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy        *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy        *User             `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentTemplateDocument) TableName() string {
	return "document_template_documents"
}
