package models

import (
	"time"
)

type ProjectChangeRequestWorkflowHistory struct {
	ID                             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectChangeRequestID         uint      `json:"project_change_request_id" gorm:"not null"`
	ProjectChangeRequestWorkflowID uint      `json:"project_change_request_workflow_id" gorm:"not null"`
	WorkflowStatusID               uint      `json:"workflow_status_id" gorm:"not null"`
	Level                          int       `json:"level" gorm:"not null"`
	Comment                        *string   `json:"comment" gorm:"type:nvarchar(max);null"`
	CreatedBYID                    uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                      time.Time `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`

	ProjectChangeRequest         *ProjectChangeRequest         `gorm:"foreignKey:ProjectChangeRequestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectChangeRequestWorkflow *ProjectChangeRequestWorkflow `gorm:"foreignKey:ProjectChangeRequestWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus               *WorkflowStatus               `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                    *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectChangeRequestWorkflowHistory) TableName() string {
	return "project_change_request_workflow_histories"
}
