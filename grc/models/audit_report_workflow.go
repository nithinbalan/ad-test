package models

import (
	"time"
)

// AuditReportWorkflow represents the structure of the audit_report_workflows table
type AuditReportWorkflow struct {
	ID                  uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditDraftReportID  *uint     `json:"audit_draft_report_id" gorm:"index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditFinalReportID  *uint     `json:"audit_final_report_id" gorm:"index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level               int       `json:"level" binding:"required" gorm:"not null"`
	ApprovalType        string    `json:"approval_type" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_type IN ('anyone', 'all')"`
	Type                string    `json:"type" binding:"required,oneof='draft' 'final'" gorm:"type:varchar(255);not null;default:'draft';check:type IN ('draft', 'final')"`
	WorkflowLevelTypeID *uint     `json:"workflow_level_type_id" gorm:"default:null"`
	WorkflowStatusID    uint      `json:"workflow_status_id" binding:"required" gorm:"not null"`
	CreatedBYID         uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relations
	AuditDraftReport  *AuditDraftReport  `gorm:"foreignKey:AuditDraftReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditFinalReport  *AuditFinalReport  `gorm:"foreignKey:AuditFinalReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowLevelType *WorkflowLevelType `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus    *WorkflowStatus    `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	AuditReportWorkflowParticipants []*AuditReportWorkflowParticipant `gorm:"foreignKey:AuditReportWorkflowID"`
	AuditReportReviews              []*AuditReportReview              `gorm:"foreignKey:AuditReportWorkflowID"`
}

// TableName overrides the default table name used by GORM
func (AuditReportWorkflow) TableName() string {
	return "audit_report_workflows"
}
