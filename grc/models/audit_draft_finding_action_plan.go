package models

import (
	"time"
)

type AuditDraftFindingActionPlan struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement"`
	AuditDraftFindingID uint       `json:"audit_draft_finding_id" binding:"required" gorm:"not null"`
	Title               string     `gorm:"type:nvarchar(500);not null"`
	Description         *string    `gorm:"type:nvarchar(MAX);null"`
	ResponsibleUserID   *uint      `json:"responsible_user_id"`
	TargetDate          *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`

	AuditDraftFinding AuditDraftFinding `gorm:"foreignKey:AuditDraftFindingID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ResponsibleUser   *User             `gorm:"foreignKey:ResponsibleUserID"`
	CreatedByUser     User              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser     *User             `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditDraftFindingActionPlan to `audit_draft_finding_action_plans`
func (AuditDraftFindingActionPlan) TableName() string {
	return "audit_draft_finding_action_plans"
}
