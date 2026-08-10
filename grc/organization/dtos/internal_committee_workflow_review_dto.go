package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	"time"
)

type InternalCommitteeWorkflowListDTO struct {
	ID                         uint                                   `json:"id"`
	InternalCommitteeID        uint                                   `json:"internal_committee_id"`
	InternalCommitteeVersionID uint                                   `json:"internal_committee_version_id"`
	NextReviewLevel            *int                                   `json:"next_review_level"`
	SubmittedAt                *time.Time                             `json:"submitted_at"`
	SubmittedByUser            *UserMinimalDTO                        `json:"submitted_by"`
	Workflows                  []*WorkflowListDTO                     `json:"workflows"`
	History                    []*InternalCommitteeWorkflowHistoryDTO `json:"history"`
}

type WorkflowListDTO struct {
	Level             int                          `json:"level"`
	ApprovalType      string                       `json:"approval_type"`
	IsSecondLastLevel bool                         `json:"is_second_last_level"`
	ReviewUsers       []*UserMinimalDTO            `json:"review_users"`
	ApprovedUsers     []*UserMinimalDTO            `json:"approved_users"`
	ApprovedAt        time.Time                    `json:"approved_at,omitempty"`
	RejectedUser      *UserMinimalDTO              `json:"rejected_user"`
	RejectedAt        time.Time                    `json:"rejected_at"`
	CurrentStatus     masterDTOs.WorkflowStatusDTO `json:"current_status"`
}

type InternalCommitteeWorkflowHistoryDTO struct {
	ID             uint                         `json:"id"`
	Level          int                          `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO `json:"workflow_status"`
	CreatedByUser  *UserMinimalDTO              `json:"created_by"`
	CreatedAt      time.Time                    `json:"created_at"`
}
