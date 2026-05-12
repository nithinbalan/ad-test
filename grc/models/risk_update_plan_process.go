package models

import "time"

type RiskUpdatePlanProcess struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint      `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	ProcessID                uint      `json:"process_id" binding:"required" gorm:"not null"`
	CreatedBYID              *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Process                *Process                `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskUpdatePlanProcess to `risk_update_plan_processes`
func (RiskUpdatePlanProcess) TableName() string {
	return "risk_update_plan_processes"
}
