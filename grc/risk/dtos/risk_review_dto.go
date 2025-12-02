package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskWorkflowListDTO struct {
	ID              uint                             `json:"id"`
	RiskRegisterID  uint                             `json:"risk_register_id"`
	SubmittedAt     *time.Time                       `json:"submitted_at"`
	SubmittedByUser *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	RiskWorkflows   []*WorkflowListDTO               `json:"workflows"`
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
