package models

import (
	"time"
)

type RiskAssociatedDocument struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	DocumentID     uint      `json:"document_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Risk          *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document     `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (RiskAssociatedDocument) TableName() string {
	return "risk_associated_documents"
}
