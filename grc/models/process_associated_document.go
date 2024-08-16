package models

import (
	"time"
)

type ProcessAssociatedDocument struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID   uint      `json:"process_id" binding:"required"`
	DocumentID  uint      `json:"document_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process       *Process  `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID"`
}

func (ProcessAssociatedDocument) TableName() string {
	return "process_associated_documents"
}
