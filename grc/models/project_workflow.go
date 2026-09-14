package models

import (
	"time"
)

type ProjectWorkflow struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID           uint       `json:"project_id" gorm:"not null"`
	Level               int        `json:"level" gorm:"not null"`
	WorkflowLevelTypeID uint       `json:"workflow_level_type_id" gorm:"not null"`
	ApprovalTypes       string     `json:"approval_types" gorm:"type:nvarchar(50);not null"`
	WorkflowStatusID    uint       `json:"workflow_status_id" gorm:"not null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime2;null"`

	Project                     *Project                      `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowLevelType           *WorkflowLevelType            `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus              *WorkflowStatus               `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                   *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectWorkflowParticipants []*ProjectWorkflowParticipant `gorm:"foreignKey:ProjectWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectWorkflowReviews      []*ProjectWorkflowReview      `gorm:"foreignKey:ProjectWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectWorkflow) TableName() string {
	return "project_workflows"
}
