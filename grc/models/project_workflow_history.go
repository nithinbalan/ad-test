package models

import (
	"time"
)

type ProjectWorkflowHistory struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID         uint      `json:"project_id" gorm:"not null"`
	ProjectWorkflowID uint      `json:"project_workflow_id" gorm:"not null"`
	WorkflowStatusID  uint      `json:"workflow_status_id" gorm:"not null"`
	Level             int       `json:"level" gorm:"not null"`
	Comment           *string   `json:"comment" gorm:"type:nvarchar(max);null"`
	CreatedBYID       uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`

	Project         *Project         `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectWorkflow *ProjectWorkflow `gorm:"foreignKey:ProjectWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus  *WorkflowStatus  `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy       *User            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectWorkflowHistory) TableName() string {
	return "project_workflow_histories"
}
