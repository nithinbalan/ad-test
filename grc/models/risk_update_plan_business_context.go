package models

import "time"

type RiskUpdatePlanBusinessContext struct {
	ID                         uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID   uint      `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	InternalAuditNote          *string   `json:"internal_audit_note" gorm:"type:nvarchar(max);default:null"`
	ExternalContext            *string   `json:"external_context" gorm:"type:nvarchar(max);default:null"`
	AnnualObjectiveDescription *string   `json:"annual_objective_description" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID                *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                  time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskUpdatePlanBusinessContext) TableName() string {
	return "risk_update_plan_business_contexts"
}
