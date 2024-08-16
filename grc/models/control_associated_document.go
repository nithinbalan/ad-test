package models

import (
	"time"
)

type ControlAssociatedDocument struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ControlID   uint      `json:"control_id" binding:"required"`
	DocumentID  uint      `json:"document_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Control       *Control  `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID"`
}

func (ControlAssociatedDocument) TableName() string {
	return "control_associated_documents"
}
