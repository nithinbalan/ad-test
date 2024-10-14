package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentVersionWorkflowListDTO struct {
	ID                uint                             `json:"id"`
	DocumentID        uint                             `json:"document_id"`
	SubmittedAt       *time.Time                       `json:"submitted_at"`
	SubmittedByUser   *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	DocumentWorkflows []*DocumentWorkflowListDTO       `json:"document_workflows"`
}

type DocumentWorkflowListDTO struct {
	Level         int                                `json:"level"`
	ApprovalType  string                             `json:"approval_type"`
	ReviewUsers   []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt    time.Time                          `json:"approved_at,omitempty"`
	CurrentStatus DocumentStatusDTO                  `json:"current_status"`
}

type DocumentReviewCommentFormDTO struct {
	Comment    string `json:"comment"`
	IsApproved bool   `json:"is_approved"`
	IsReverted bool   `json:"is_reverted"`
	IsRejected bool   `json:"is_rejected"`
}
