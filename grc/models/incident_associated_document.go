package models

import (
	"time"
)

type IncidentAssociatedDocument struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	DocumentID         uint      `json:"document_id" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document         *Document         `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
}

func (IncidentAssociatedDocument) TableName() string {
	return "incident_associated_documents"
}
