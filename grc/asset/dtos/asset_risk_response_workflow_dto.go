package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AssetRiskResponseWorkflowListDTO struct {
	ID                    uint                                   `json:"id"`
	RiskResponseID        uint                                   `json:"risk_response_id"`
	NextReviewLevel       *int                                   `json:"next_review_level"`
	SubmittedAt           *time.Time                             `json:"submitted_at"`
	SubmittedByUser       *organizationDTOs.UserMinimalDTO       `json:"submitted_by"`
	ApprovedAt            *time.Time                             `json:"approved_at"`
	ApprovedByUser        *organizationDTOs.UserMinimalDTO       `json:"approved_by"`
	CurrentStatus         *masterDTOs.WorkflowStatusDTO          `json:"workflow_status"`
	RiskResponseWorkflows []*WorkflowListDTO                     `json:"workflows"`
	History               []*AssetRiskResponseWorkflowHistoryDTO `json:"history"`
}

type AssetRiskResponseWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}
