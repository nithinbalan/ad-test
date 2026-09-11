package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectWorkflowListDTO struct {
	ID              uint                             `json:"id"`
	ProjectID       uint                             `json:"project_id"`
	SubmittedAt     *time.Time                       `json:"submitted_at"`
	SubmittedByUser *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	ApprovedAt      *time.Time                       `json:"approved_at"`
	ApprovedByUser  *organizationDTOs.UserMinimalDTO `json:"approved_by"`
	RejectedAt      *time.Time                       `json:"rejected_at"`
	RejectedByUser  *organizationDTOs.UserMinimalDTO `json:"rejected_by"`
	NextReviewLevel *int                             `json:"next_review_level"`
	WorkflowStatus  *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	Workflows       []*ProjectWorkflowLevelDTO       `json:"workflows"`
	History         []*ProjectWorkflowHistoryDTO     `json:"history"`
}

type ProjectWorkflowLevelDTO struct {
	Level            int                                `json:"level"`
	ApprovalType     string                             `json:"approval_type"`
	ReviewUsers      []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers    []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt       *time.Time                         `json:"approved_at"`
	ApprovedAtString string                             `json:"approved_at_string"`
	RejectedUser     *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	RejectedAt       *time.Time                         `json:"rejected_at"`
	CurrentStatus    masterDTOs.WorkflowStatusDTO       `json:"current_status"`
}

type ProjectWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	Comment        *string                          `json:"comment"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}

type ProjectReviewFormDTO struct {
	Comment *string `json:"comment"`
}

type ProjectWorkflowReviewFormDTO = ProjectReviewFormDTO
