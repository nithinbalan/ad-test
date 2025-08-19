package models

import (
	"time"
)

type KRIAssociatedDocument struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRIRegisterID uint      `json:"kri_register_id" binding:"required"`
	DocumentID    uint      `json:"document_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KRIRegister   *KRIRegister `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document    `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (KRIAssociatedDocument) TableName() string {
	return "kri_associated_documents"
}
