package models

import (
	"time"
)

// InternalCommitteeWorkflow represents the structure of internal committee workflows
type InternalCommitteeWorkflow struct {
	ID                         uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	InternalCommitteeID        uint       `json:"internal_committee_id" binding:"required" gorm:"not null;index"`
	InternalCommitteeVersionID uint       `json:"internal_committee_version_id" binding:"required" gorm:"not null;index"`
	Level                      int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowLevelTypeID        uint       `json:"workflow_level_type_id" binding:"required" gorm:"not null"`
	ApprovalTypes              string     `json:"approval_types" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_types IN ('anyone','all')"`
	WorkflowStatusID           uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	CreatedBYID                uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	//Relations
	InternalCommittee        *InternalCommittee        `gorm:"foreignKey:InternalCommitteeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InternalCommitteeVersion *InternalCommitteeVersion `gorm:"foreignKey:InternalCommitteeVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowLevelType        *WorkflowLevelType        `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus           *WorkflowStatus           `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                *User                     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                *User                     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	InternalCommitteeWorkflowParticipants []*InternalCommitteeWorkflowParticipant `gorm:"foreignKey:InternalCommitteeWorkflowID"`
	InternalCommitteeWorkflowReviews      []*InternalCommitteeWorkflowReview      `gorm:"foreignKey:InternalCommitteeWorkflowID"`
}

// TableName overrides table name
func (InternalCommitteeWorkflow) TableName() string {
	return "internal_committee_workflows"
}
