package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IdentityReviewCampaignWorkflowListDTO struct {
	ID              uint                                        `json:"id"`
	CampaignID      uint                                        `json:"identity_review_campaign_id"`
	SubmittedAt     *time.Time                                  `json:"submitted_at"`
	SubmittedByUser *organizationDTOs.UserMinimalDTO            `json:"submitted_by"`
	ApprovedAt      *time.Time                                  `json:"approved_at"`
	ApprovedByUser  *organizationDTOs.UserMinimalDTO            `json:"approved_by"`
	NextReviewLevel *int                                        `json:"next_review_level"`
	CurrentStatus   *masterDTOs.WorkflowStatusDTO               `json:"workflow_status"`
	Workflows       []*IdentityReviewCampaignWorkflowLevelDTO   `json:"workflows"`
	History         []*IdentityReviewCampaignWorkflowHistoryDTO `json:"history"`
}

type IdentityReviewCampaignWorkflowLevelDTO struct {
	Level                 int                                `json:"level"`
	ApprovalType          string                             `json:"approval_type"`
	ReviewUsers           []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers         []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt            *time.Time                         `json:"approved_at"`
	ApprovedAtString      string                             `json:"approved_at_string"`
	RejectedUser          *organizationDTOs.UserMinimalDTO   `json:"rejected_user"`
	RejectedAt            *time.Time                         `json:"rejected_at"`
	CurrentStatus         masterDTOs.WorkflowStatusDTO       `json:"current_status"`
	StatusTextColor       string                             `json:"-"`
	StatusBorderColor     string                             `json:"-"`
	StatusBackgroundColor string                             `json:"-"`
}

type IdentityReviewCampaignWorkflowHistoryDTO struct {
	ID             uint                             `json:"id"`
	Level          int                              `json:"level"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	CreatedByUser  *organizationDTOs.UserMinimalDTO `json:"created_by"`
	CreatedAt      time.Time                        `json:"created_at"`
}
