package models

import (
	"time"
)

type AuditFindingActionPlanStatusUpdate struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingActionPlanID uint       `json:"audit_finding_action_plan_id" gorm:"not null"`
	ActionPlanStatusID       uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage               *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                  string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditFindingActionPlan *AuditFindingActionPlan                       `gorm:"foreignKey:AuditFindingActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus       *ActionPlanStatus                             `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                                         `gorm:"foreignKey:UpdatedByID"`
	Documents              []*AuditFindingActionPlanStatusUpdateDocument `gorm:"foreignKey:AuditFindingActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by AuditFindingActionPlanStatusUpdate to `audit_finding_action_plan_status_updates`
func (AuditFindingActionPlanStatusUpdate) TableName() string {
	return "audit_finding_action_plan_status_updates"
}
