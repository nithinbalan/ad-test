package models

import "time"

type ProjectSecurityGateChecklistExceptionWorkflowReview struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ExceptionWorkflowID uint       `json:"exception_workflow_id" gorm:"not null"`
	WorkflowStatusID    uint       `json:"workflow_status_id" gorm:"not null"`
	Comment             *string    `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ExceptionWorkflow *ProjectSecurityGateChecklistExceptionWorkflow `json:"exception_workflow" gorm:"foreignKey:ExceptionWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus    *WorkflowStatus                                `json:"workflow_status" gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                                          `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser     *User                                          `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistExceptionWorkflowReview) TableName() string {
	return "project_security_gate_checklist_exception_workflow_reviews"
}
