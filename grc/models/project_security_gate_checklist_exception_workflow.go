package models

import "time"

type ProjectSecurityGateChecklistExceptionWorkflow struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ChecklistID         uint       `json:"checklist_id" gorm:"not null"`
	Level               int        `json:"level" gorm:"not null"`
	WorkflowLevelTypeID uint       `json:"workflow_level_type_id" gorm:"not null"`
	ApprovalTypes       string     `json:"approval_types" gorm:"type:nvarchar(50);not null"`
	WorkflowStatusID    uint       `json:"workflow_status_id" gorm:"not null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Checklist         *ProjectSecurityGateChecklist `json:"checklist" gorm:"foreignKey:ChecklistID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowLevelType *WorkflowLevelType            `json:"workflow_level_type" gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus    *WorkflowStatus               `json:"workflow_status" gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser     *User                         `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	Participants []*ProjectSecurityGateChecklistExceptionWorkflowParticipant `json:"participants" gorm:"foreignKey:ExceptionWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Reviews      []*ProjectSecurityGateChecklistExceptionWorkflowReview      `json:"reviews" gorm:"foreignKey:ExceptionWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistExceptionWorkflow) TableName() string {
	return "project_security_gate_checklist_exception_workflows"
}
