package models

import (
	"time"
)

// DocumentChangeRequestWorkflow represents the structure of the document change request workflow
type DocumentChangeRequestWorkflow struct {
	ID                                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // Primary key, auto-incremented
	DocumentChangeRequestID               uint       `json:"document_change_request_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level                                 int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowLevelTypeID                   uint       `json:"workflow_level_type_id" binding:"required" gorm:"not null"`
	ApprovalTypes                         string     `json:"approval_types" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_types IN ('anyone', 'all')"`
	DocumentChangeRequestWorkflowStatusID uint       `json:"document_change_request_workflow_status_id" binding:"required" gorm:"not null"`
	CreatedBYID                           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentChangeRequest               *DocumentChangeRequest               `gorm:"foreignKey:DocumentChangeRequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowLevelType                   *WorkflowLevelType                   `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestWorkflowStatus *DocumentChangeRequestWorkflowStatus `gorm:"foreignKey:DocumentChangeRequestWorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                           *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                           *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	DocumentChangeRequestWorkflowParticipants []*DocumentChangeRequestWorkflowParticipant `gorm:"foreignKey:DocumentChangeRequestWorkflowID"`
	DocumentChangeRequestWorkflowReviews      []*DocumentChangeRequestWorkflowReview      `gorm:"foreignKey:DocumentChangeRequestWorkflowID"`
}

// TableName overrides the table name used by User to `document_change_request_workflows`
func (DocumentChangeRequestWorkflow) TableName() string {
	return "document_change_request_workflows"
}
