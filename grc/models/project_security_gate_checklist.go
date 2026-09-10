package models

import "time"

type ProjectSecurityGateChecklist struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SecurityGateID         uint       `json:"security_gate_id" gorm:"not null"`
	RequirementID          uint       `json:"requirement_id" gorm:"not null"`
	ResponseStatusID       *uint      `json:"response_status_id" gorm:"default:null"`
	Comment                *string    `json:"comment" gorm:"type:nvarchar(max);default:null"`
	IsException            bool       `json:"is_exception" gorm:"not null;default:false"`
	ExceptionJustification *string    `json:"exception_justification" gorm:"type:nvarchar(max);default:null"`
	WorkflowStatusID       *uint      `json:"workflow_status_id" gorm:"default:null"`
	SubmittedBYID          *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt            *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	NextReviewLevel        *int       `json:"next_review_level" gorm:"default:null"`
	ApprovedBYID           *uint      `json:"approved_by_id" gorm:"default:null"`
	ApprovedAt             *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	CreatedBYID            *uint      `json:"created_by_id" gorm:"default:null"`
	CreatedAt              *time.Time `json:"created_at" gorm:"type:datetime;default:null"`
	UpdatedBYID            *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	SecurityGate    *ProjectSecurityGate                             `json:"security_gate" gorm:"foreignKey:SecurityGateID"`
	Requirement     *ProjectSecurityGateChecklistTemplateRequirement `json:"requirement" gorm:"foreignKey:RequirementID"`
	ResponseStatus  *ProjectSecurityGateChecklistResponseStatus      `json:"response_status" gorm:"foreignKey:ResponseStatusID"`
	WorkflowStatus  *WorkflowStatus                                  `json:"workflow_status" gorm:"foreignKey:WorkflowStatusID"`
	SubmittedByUser *User                                            `json:"submitted_by_user" gorm:"foreignKey:SubmittedBYID"`
	ApprovedByUser  *User                                            `json:"approved_by_user" gorm:"foreignKey:ApprovedBYID"`
	CreatedByUser   *User                                            `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                                            `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID"`
	Documents       []ProjectSecurityGateChecklistResponseDocument   `json:"documents" gorm:"foreignKey:ChecklistID"`
	Comments        []ProjectSecurityGateChecklistComment            `json:"comments" gorm:"foreignKey:ChecklistID"`
	Controls        []ProjectSecurityGateChecklistControl            `json:"controls" gorm:"foreignKey:ChecklistID"`

	ExceptionWorkflows []*ProjectSecurityGateChecklistExceptionWorkflow `json:"exception_workflows" gorm:"foreignKey:ChecklistID"`
}

func (ProjectSecurityGateChecklist) TableName() string {
	return "project_security_gate_checklists"
}
