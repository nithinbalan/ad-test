package models

import (
	"time"
)

// AuditReportWorkflowParticipant represents a participant in a audit_report workflow.
type AuditReportWorkflowParticipant struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditReportWorkflowID uint       `json:"audit_report_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID                uint       `json:"user_id" binding:"required" gorm:"not null"`
	CreatedBYID           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relations
	AuditReportWorkflow *AuditReportWorkflow `gorm:"foreignKey:AuditReportWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                *User                `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy           *User                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by AuditReportWorkflowParticipant to `audit_report_workflow_participants`
func (AuditReportWorkflowParticipant) TableName() string {
	return "audit_report_workflow_participants"
}
