package models

import (
	"time"
)

// DocumentChangeRequestWorkflowReview represents the structure of our resource
type DocumentChangeRequestWorkflowReview struct {
	ID                                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // Primary key, auto-incremented
	DocumentChangeRequestWorkflowID       uint       `json:"document_change_request_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DocumentChangeRequestWorkflowStatusID uint       `json:"document_change_request_workflow_status_id" binding:"required" gorm:"not null"`
	Comment                               string     `json:"comment" binding:"-" gorm:"type:text"`
	CreatedBYID                           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentChangeRequestWorkflow       *DocumentChangeRequestWorkflow       `gorm:"foreignKey:DocumentChangeRequestWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestWorkflowStatus *DocumentChangeRequestWorkflowStatus `gorm:"foreignKey:DocumentChangeRequestWorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                       *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser                       *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by DocumentChangeRequestWorkflowReview to `document_change_request_workflow_reviews`
func (DocumentChangeRequestWorkflowReview) TableName() string {
	return "document_change_request_workflow_reviews"
}
