package models

import (
	"time"
)

type KriIncident struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KriRegisterID      uint      `json:"kri_register_id" binding:"required"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KriRegister      *KRIRegister      `gorm:"foreignKey:KriRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
}

func (KriIncident) TableName() string {
	return "kri_incidents"
}
