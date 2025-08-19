package models

import (
	"time"
)

type ProjectAssociatedDocument struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID   uint      `json:"project_id" binding:"required"`
	DocumentID  uint      `json:"document_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Project       *Project  `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID"`
}

func (ProjectAssociatedDocument) TableName() string {
	return "project_associated_documents"
}
