package models

import (
	"time"
)

type BCPControl struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" binding:"required"`
	ControlID                uint      `json:"control_id" binding:"required"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control                *Control                `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
}

func (BCPControl) TableName() string {
	return "bcp_controls"
}
