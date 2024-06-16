package models

import (
	"time"
)

type IncidentDocument struct {
	ID                 uint      `json:"id" gorm:"primary_key;auto_increment"`
	IncidentRegisterID uint      `json:"incident_register_id"`
	Language           string    `json:"language" binding:"required"`
	URL                string    `json:"url"`
	Name               string    `json:"name" binding:"required"`
	Size               int       `json:"size"`
	Token              string    `json:"token"`
	OriginalName       string    `json:"original_name"`
	Type               string    `json:"type"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBy          uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBy"`
}

func (IncidentDocument) TableName() string {
	return "incident_documents"
}
