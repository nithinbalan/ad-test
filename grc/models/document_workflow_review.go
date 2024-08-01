package models

import (
	"time"
)

// DocumentWorkflowReview represents the structure of our resource
type DocumentWorkflowReview struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // Primary key, auto-incremented
	DocumentWorkflowID       uint       `json:"document_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DocumentWorkflowStatusID uint       `json:"document_workflow_status_id" binding:"required" gorm:"not null"`
	Comment                  string     `json:"comment" binding:"-" gorm:"type:text"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentWorkflow       *DocumentWorkflow       `gorm:"foreignKey:DocumentWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentWorkflowStatus *DocumentWorkflowStatus `gorm:"foreignKey:DocumentWorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by DocumentWorkflowReview to `document_workflow_reviews`
func (DocumentWorkflowReview) TableName() string {
	return "document_workflow_reviews"
}
