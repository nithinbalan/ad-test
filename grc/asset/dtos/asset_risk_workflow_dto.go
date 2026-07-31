package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AssetRiskWorkflowListDTO struct {
	ID              uint                             `json:"id"`
	RiskRegisterID  uint                             `json:"risk_register_id"`
	SubmittedAt     *time.Time                       `json:"submitted_at"`
	SubmittedByUser *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	ApprovedAt      *time.Time                       `json:"approved_at"`
	ApprovedByUser  *organizationDTOs.UserMinimalDTO `json:"approved_by"`
	NextReviewLevel *int                             `json:"next_review_level"`
	CurrentStatus   *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	RiskWorkflows   []*WorkflowListDTO               `json:"workflows"`
	History         []*AssetRiskWorkflowHistoryDTO   `json:"history"`
}

type WorkflowListDTO struct {
	Level         int                                `json:"level"`
	ApprovalType  string                             `json:"approval_type"`
	ReviewUsers   []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt    *time.Time                         `json:"approved_at"`
	RejectedUser  *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	RejectedAt    *time.Time                         `json:"rejected_at"`
	CurrentStatus masterDTOs.WorkflowStatusDTO       `json:"current_status"`
}

type AssetRiskWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}
