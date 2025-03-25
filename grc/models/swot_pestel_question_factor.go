package models

import (
	"time"
)

// AuditReportWorkflowParticipant represents a participant in a audit_report workflow.
type SwotPestelQuestionFactor struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SwotPestelQuestionID uint       `json:"swot_pestel_question_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SwotPestelCategoryID uint       `json:"swot_pestel_category_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Description          string     `json:"description" binding:"max=250" gorm:"type:varchar(max)"`
	CreatedBYID          uint       `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID          *uint      `json:"updated_by_id"`
	UpdatedAt            *time.Time `json:"updated_at"`

	// Relations
	SwotPestelQuestion *SwotPestelQuestion `gorm:"foreignKey:SwotPestelQuestionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SwotPestelCategory *SwotPestelCategory `gorm:"foreignKey:SwotPestelCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditReportWorkflowParticipant to `audit_report_workflow_participants`
func (SwotPestelQuestionFactor) TableName() string {
	return "swot_pestel_question_factors"
}
