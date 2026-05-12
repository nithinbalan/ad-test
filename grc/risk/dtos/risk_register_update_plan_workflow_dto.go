package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskRegisterUpdatePlanWorkflowListDTO struct {
	ID                              uint                                        `json:"id"`
	RiskRegisterUpdatePlanID        uint                                        `json:"risk_register_update_plan_id"`
	NextReviewLevel                 *int                                        `json:"next_review_level"`
	SubmittedAt                     *time.Time                                  `json:"submitted_at"`
	SubmittedByUser                 *organizationDTOs.UserMinimalDTO            `json:"submitted_by"`
	RiskRegisterUpdatePlanWorkflows []*WorkflowListDTO                          `json:"workflows"`
	History                         []*RiskRegisterUpdatePlanWorkflowHistoryDTO `json:"history"`
}

type RiskRegisterUpdatePlanWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
	Comment        string                           `json:"comment"`
}
