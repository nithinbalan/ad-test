package models

import (
	"time"
)

type AuditFindingActionPlanComment struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoKncrement"`
	AuditFindingActionPlanID uint      `json:"audit_finding_action_plan_id" gorm:"not null"`
	Comment                  string    `json:"comment" binding:"required"`
	CreatedBYID              uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFindingActionPlan *AuditFindingActionPlan                  `gorm:"foreignKey:AuditFindingActionPlanID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	CreatedByUser          *User                                    `gorm:"foreignKey:CreatedBYID"`
	Documents              *[]AuditFindingActionPlanCommentDocument `gorm:"foreignKey:AuditFindingActionPlanCommentID"`
}

func (AuditFindingActionPlanComment) TableName() string {
	return "audit_finding_action_plan_comments"
}
