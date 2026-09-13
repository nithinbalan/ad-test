package dtos

import (
	"time"

	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectSecurityGateChecklistExceptionWorkflowListDTO struct {
	ID              uint                                                           `json:"id"`
	ChecklistID     uint                                                           `json:"checklist_id"`
	NextReviewLevel *int                                                           `json:"next_review_level"`
	SubmittedAt     *time.Time                                                     `json:"submitted_at"`
	SubmittedByUser *organizationDTOs.UserMinimalDTO                               `json:"submitted_by_user"`
	ApprovedAt      *time.Time                                                     `json:"approved_at"`
	ApprovedByUser  *organizationDTOs.UserMinimalDTO                               `json:"approved_by_user"`
	WorkflowStatus  *masterDTOs.WorkflowStatusDTO                                  `json:"workflow_status"`
	Workflows       []*ProjectSecurityGateChecklistExceptionWorkflowLevelDTO       `json:"workflows"`
	History         []*ProjectSecurityGateChecklistExceptionWorkflowHistoryItemDTO `json:"history"`
}

type ProjectSecurityGateChecklistExceptionWorkflowLevelDTO struct {
	Level          int                                `json:"level"`
	ApprovalType   string                             `json:"approval_type"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO       `json:"workflow_status"`
	ReviewUsers    []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers  []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	RevertedUser   *organizationDTOs.UserMinimalDTO   `json:"reverted_user"`
	ApprovedAt     *time.Time                         `json:"approved_at"`
	RevertedAt     *time.Time                         `json:"reverted_at"`
}

type ProjectSecurityGateChecklistExceptionWorkflowHistoryItemDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	Comment        *string                          `json:"comment"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	CreatedAt      time.Time                        `json:"created_at"`
}
