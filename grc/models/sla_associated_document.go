package models

import (
	"time"
)

type SLAAssociatedDocument struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLARegisterID uint      `json:"sla_register_id" binding:"required"`
	DocumentID    uint      `json:"document_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SLARegister   *SLARegister `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document    `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (SLAAssociatedDocument) TableName() string {
	return "sla_associated_documents"
}
