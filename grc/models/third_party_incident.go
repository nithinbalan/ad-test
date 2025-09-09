package models

import (
	"time"
)

type ThirdPartyIncident struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	IncidentRegisterID   uint      `json:"incident_register_id" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister   *IncidentRegister   `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

func (ThirdPartyIncident) TableName() string {
	return "third_party_incidents"
}
