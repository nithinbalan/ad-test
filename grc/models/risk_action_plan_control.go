package models

import (
	"time"
)

type RiskActionPlanControl struct {
	ID               uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskActionPlanID uint      `json:"risk_action_plan_id" binding:"required"`
	ControlID        uint      `json:"control_id" binding:"required"`
	CreatedByID      uint      `json:"created_by_id" binding:"required"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskActionPlan *RiskActionPlan `gorm:"foreignKey:RiskActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control        *Control        `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID"`
}

func (RiskActionPlanControl) TableName() string {
	return "risk_action_plan_controls"
}
