package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BiaProcessWorkflowListDTO struct {
	ID                  uint                             `json:"id"`
	BiaProcessID        uint                             `json:"bia_process_id"`
	NextReviewLevel     *int                             `json:"next_review_level"`
	SubmittedAt         *time.Time                       `json:"submitted_at"`
	SubmittedByUser     *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	BiaProcessWorkflows []*BiaWorkflowListDTO            `json:"workflows"`
	History             []*BiaProcessWorkflowHistoryDTO  `json:"history"`
}

type BiaProcessWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}

type BiaWorkflowListDTO struct {
	Level         int                                `json:"level"`
	ApprovalType  string                             `json:"approval_type"`
	ReviewUsers   []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt    time.Time                          `json:"approved_at,omitempty"`
	RejectedUser  *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	RejectedAt    time.Time                          `json:"rejected_at"`
	CurrentStatus masterDTOs.WorkflowStatusDTO       `json:"current_status"`
}
