package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type KRIRegisterWorkflowListDTO struct {
	ID                   uint                             `json:"id"`
	KRIRegisterID        uint                             `json:"kri_register_id"`
	SubmittedAt          *time.Time                       `json:"submitted_at"`
	SubmittedByUser      *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	KRIRegisterWorkflows []*WorkflowListDTO               `json:"workflows"`
}
