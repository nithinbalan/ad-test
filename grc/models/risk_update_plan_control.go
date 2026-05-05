package models

import "time"

type RiskUpdatePlanControl struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint      `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	ControlID                uint      `json:"control_id" binding:"required" gorm:"not null"`
	CreatedBYID              *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Control                *Control                `gorm:"foreignKey:ControlID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskUpdatePlanControl to `risk_update_plan_controls`
func (RiskUpdatePlanControl) TableName() string {
	return "risk_update_plan_controls"
}
