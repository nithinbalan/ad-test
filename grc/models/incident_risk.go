package models

import (
	"time"
)

type IncidentRisk struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	RiskRegisterID     uint      `json:"risk_register_id" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister     *RiskRegister     `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
}

func (IncidentRisk) TableName() string {
	return "incident_risks"
}
