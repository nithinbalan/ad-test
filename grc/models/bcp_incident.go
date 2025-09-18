package models

import (
	"time"
)

type BcpIncident struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" binding:"required"`
	IncidentRegisterID       uint      `json:"incident_register_id" binding:"required"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister       *IncidentRegister       `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
}

func (BcpIncident) TableName() string {
	return "bcp_incidents"
}
