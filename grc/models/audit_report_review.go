package models

import (
	"time"
)

// AuditReportReview represents the structure of the audit_report_reviews table
type AuditReportReview struct {
	ID                    uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditDraftReportID    *uint     `json:"audit_draft_report_id" gorm:"default:null"`
	AuditReportWorkflowID uint      `json:"audit_report_workflow_id" binding:"required" gorm:"not null"`
	WorkflowStatusID      uint      `json:"workflow_status_id" binding:"required" gorm:"not null"`
	UserID                uint      `json:"user_id" binding:"required" gorm:"not null"`
	Comment               *string   `json:"comment" gorm:"type:nvarchar(MAX);default:null"`
	CreatedBYID           uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt             time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relations
	AuditDraftReport    *AuditDraftReport    `gorm:"foreignKey:AuditDraftReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditReportWorkflow *AuditReportWorkflow `gorm:"foreignKey:AuditReportWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus      *WorkflowStatus      `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                *User                `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the default table name used by GORM
func (AuditReportReview) TableName() string {
	return "audit_report_reviews"
}
