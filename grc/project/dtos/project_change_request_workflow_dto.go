package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectChangeRequestWorkflowListDTO struct {
	ID                     uint                                      `json:"id"`
	ProjectChangeRequestID uint                                      `json:"project_change_request_id"`
	SubmittedAt            *time.Time                                `json:"submitted_at"`
	SubmittedByUser        *organizationDTOs.UserMinimalDTO          `json:"submitted_by_user"`
	ApprovedAt             *time.Time                                `json:"approved_at"`
	ApprovedByUser         *organizationDTOs.UserMinimalDTO          `json:"approved_by_user"`
	RejectedAt             *time.Time                                `json:"rejected_at"`
	RejectedByUser         *organizationDTOs.UserMinimalDTO          `json:"rejected_by_user"`
	NextReviewLevel        *int                                      `json:"next_review_level"`
	WorkflowStatus         *masterDTOs.WorkflowStatusDTO             `json:"workflow_status"`
	Workflows              []*ProjectChangeRequestWorkflowLevelDTO   `json:"workflows"`
	History                []*ProjectChangeRequestWorkflowHistoryDTO `json:"history"`
}

type ProjectChangeRequestWorkflowLevelDTO struct {
	Level            int                                `json:"level"`
	ApprovalType     string                             `json:"approval_type"`
	WorkflowStatus   masterDTOs.WorkflowStatusDTO       `json:"workflow_status"`
	ReviewUsers      []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers    []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	RejectedUser     *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	ApprovedAt       *time.Time                         `json:"approved_at"`
	ApprovedAtString string                             `json:"approved_at_string"`
	RejectedAt       *time.Time                         `json:"rejected_at"`
}

type ProjectChangeRequestWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	Comment        *string                          `json:"comment"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	CreatedAt      time.Time                        `json:"created_at"`
}

type ProjectChangeRequestReviewFormDTO struct {
	Comment *string `json:"comment"`
}
