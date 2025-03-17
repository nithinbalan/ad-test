package models

import (
	"time"
)

// AuditReportWorkflowParticipant represents a participant in a audit_report workflow.
type SwotPestelQuestionInvitee struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SwotPestelQuestionID uint      `json:"swot_pestel_question_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID               uint      `json:"user_id" binding:"required" gorm:"not null"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relations
	SwotPestelQuestion *SwotPestelQuestion `gorm:"foreignKey:SwotPestelQuestionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User               *User               `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy          *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by AuditReportWorkflowParticipant to `audit_report_workflow_participants`
func (SwotPestelQuestionInvitee) TableName() string {
	return "swot_pestel_question_invitees"
}
