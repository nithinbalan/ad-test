package models

import (
	"time"
)

type ProjectDecisionLogWorkflowReview struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectDecisionLogWorkflowID uint       `json:"project_decision_log_workflow_id" gorm:"not null"`
	WorkflowStatusID             uint       `json:"workflow_status_id" gorm:"not null"`
	Comment                      *string    `json:"comment" gorm:"type:nvarchar(max);null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime2;null"`

	ProjectDecisionLogWorkflow *ProjectDecisionLogWorkflow `gorm:"foreignKey:ProjectDecisionLogWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus             *WorkflowStatus             `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                  *User                       `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectDecisionLogWorkflowReview) TableName() string {
	return "project_decision_log_workflow_reviews"
}
