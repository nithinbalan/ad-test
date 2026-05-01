package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskResponseWorkflowListDTO struct {
	ID                   uint                              `json:"id"`
	RiskResponseID       uint                              `json:"risk_response_id"`
	NextReviewLevel      *int                              `json:"next_review_level"`
	SubmittedAt          *time.Time                        `json:"submitted_at"`
	SubmittedByUser      *organizationDTOs.UserMinimalDTO  `json:"submitted_by"`
	RiskResponseWorkflows []*WorkflowListDTO               `json:"workflows"`
	History              []*RiskResponseWorkflowHistoryDTO  `json:"history"`
}

type RiskResponseWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}
