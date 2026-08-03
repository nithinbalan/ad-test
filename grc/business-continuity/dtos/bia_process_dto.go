package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BiaProcessDetailDTO struct {
	ID              uint                             `json:"id"`
	ProcessID       uint                             `json:"process_id"`
	WorkflowStatus  *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	NextReviewLevel *int                             `json:"next_review_level"`
	IsCyber         bool                             `json:"is_cyber"`
	SubmittedAt     *time.Time                       `json:"submitted_at"`
	SubmittedBy     *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	ApprovedAt      *time.Time                       `json:"approved_at"`
	ApprovedBy      *organizationDTOs.UserMinimalDTO `json:"approved_by"`
	CreatedBYID     *uint                            `json:"created_by_id"`
	CreatedBy       *organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedBYID     *uint                            `json:"updated_by_id"`
	UpdatedBy       *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
}
