package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type KRIRegisterWorkflowListDTO struct {
	ID                   uint                             `json:"id"`
	KRIRegisterID        uint                             `json:"kri_register_id"`
	NextReviewLevel      *int                             `json:"next_review_level"`
	SubmittedAt          *time.Time                       `json:"submitted_at"`
	SubmittedByUser      *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	KRIRegisterWorkflows []*WorkflowListDTO               `json:"workflows"`
	History              []*KRIRegisterWorkflowHistoryDTO `json:"history"`
}

type KRIRegisterWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}
