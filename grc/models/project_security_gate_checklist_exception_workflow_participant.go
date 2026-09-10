package models

import "time"

type ProjectSecurityGateChecklistExceptionWorkflowParticipant struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ExceptionWorkflowID uint       `json:"exception_workflow_id" gorm:"not null"`
	RoleID              *uint      `json:"role_id" gorm:"default:null"`
	UserID              uint       `json:"user_id" gorm:"not null"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ExceptionWorkflow *ProjectSecurityGateChecklistExceptionWorkflow `json:"exception_workflow" gorm:"foreignKey:ExceptionWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Role              *Role                                          `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User              *User                                          `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                                          `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser     *User                                          `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistExceptionWorkflowParticipant) TableName() string {
	return "project_security_gate_checklist_exception_workflow_participants"
}
