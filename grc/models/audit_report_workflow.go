package models

import (
	"time"
)

// AuditReportWorkflow represents the structure of the audit_report_workflows table
type AuditReportWorkflow struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditDraftReportID    *uint      `json:"audit_draft_report_id" gorm:"index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level                 int        `json:"level" binding:"required" gorm:"not null"`
	ApprovalType          string     `json:"approval_type" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_type IN ('anyone', 'all')"`
	Type                  string     `json:"type" binding:"required,oneof='draft-report' 'final-report'" gorm:"type:varchar(255);not null;default:'draft-report';check:type IN ('draft-report', 'final-report')"`
	AuditWorkflowStatusID uint       `json:"audit_workflow_status_id" binding:"required" gorm:"not null"`
	Comment               *string    `json:"comment" gorm:"type:nvarchar(MAX);default:null"`
	CreatedBYID           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	AuditDraftReport    *AuditDraftReport    `gorm:"foreignKey:AuditDraftReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditWorkflowStatus *AuditWorkflowStatus `gorm:"foreignKey:AuditWorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy           *User                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedBy           *User                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	AuditReportWorkflowParticipants []*AuditReportWorkflowParticipant `gorm:"foreignKey:AuditReportWorkflowID"`
}

// TableName overrides the default table name used by GORM
func (AuditReportWorkflow) TableName() string {
	return "audit_report_workflows"
}
