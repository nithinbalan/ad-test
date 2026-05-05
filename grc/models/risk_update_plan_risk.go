package models

import "time"

type RiskUpdatePlanRisk struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint      `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	RiskRegisterID           uint      `json:"risk_register_id" binding:"required" gorm:"not null"`
	CreatedBYID              *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskRegister           *RiskRegister           `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskUpdatePlanRisk to `risk_update_plan_risks`
func (RiskUpdatePlanRisk) TableName() string {
	return "risk_update_plan_risks"
}
