package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IncidentClosureWorkflowListDTO struct {
	ID                       uint                                 `json:"id"`
	IncidentClosureID        uint                                 `json:"incident_closure_id"`
	NextReviewLevel          *int                                 `json:"next_review_level"`
	SubmittedAt              *time.Time                           `json:"submitted_at"`
	SubmittedByUser          *organizationDTOs.UserMinimalDTO     `json:"submitted_by"`
	IncidentClosureWorkflows []*WorkflowListDTO                   `json:"workflows"`
	History                  []*IncidentClosureWorkflowHistoryDTO `json:"history"`
}

type IncidentClosureWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
	Comment        string                           `json:"comment"`
}

type WorkflowListDTO struct {
	Level         int                                `json:"level"`
	ApprovalType  string                             `json:"approval_type"`
	ReviewUsers   []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt    time.Time                          `json:"approved_at,omitempty"`
	RejectedUser  *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	RejectedAt    time.Time                          `json:"rejected_at"`
	CurrentStatus masterDTOs.WorkflowStatusDTO       `json:"current_status"`
}
