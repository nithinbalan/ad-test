package models

import (
	"time"
)

// InternalCommitteeWorkflowReview represents workflow review for internal committee
type InternalCommitteeWorkflowReview struct {
	ID                          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	InternalCommitteeWorkflowID uint       `json:"internal_committee_workflow_id" binding:"required" gorm:"not null;index"`
	WorkflowStatusID            uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment                     string     `json:"comment" binding:"-" gorm:"type:nvarchar(MAX)"`
	CreatedBYID                 uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                   time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	//  Relations
	InternalCommitteeWorkflow *InternalCommitteeWorkflow `gorm:"foreignKey:InternalCommitteeWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus            *WorkflowStatus            `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides table name
func (InternalCommitteeWorkflowReview) TableName() string {
	return "internal_committee_workflow_reviews"
}
