package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentCRWorkflowListDTO struct {
	ID                uint                             `json:"id"`
	DocumentID        uint                             `json:"document_id"`
	SubmittedAt       *time.Time                       `json:"submitted_at"`
	SubmittedByUser   *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	DocumentWorkflows []*DocumentWorkflowListDTO       `json:"document_workflows"`
}

type DocumentCRReviewCommentFormDTO struct {
	Comment    string `json:"comment"`
	IsApproved bool   `json:"is_approved"`
	IsReverted bool   `json:"is_reverted"`
	IsRejected bool   `json:"is_rejected"`
}
