package models

import "time"

type ProjectSecurityGateChecklistExceptionWorkflowHistory struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ChecklistID      uint      `json:"checklist_id" gorm:"not null"`
	WorkflowStatusID uint      `json:"workflow_status_id" gorm:"not null"`
	Level            int       `json:"level" gorm:"not null"`
	Comment          *string   `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Checklist      *ProjectSecurityGateChecklist `json:"checklist" gorm:"foreignKey:ChecklistID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus *WorkflowStatus               `json:"workflow_status" gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistExceptionWorkflowHistory) TableName() string {
	return "project_security_gate_checklist_exception_workflow_histories"
}
