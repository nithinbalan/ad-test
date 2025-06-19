package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type GRCReportWorkflowDetailDTO struct {
	ID                      uint                             `json:"id"`
	Type                    string                           `json:"type"`
	GRCReportID             uint                             `json:"grc_report_id"`
	GRCReport               GRCReportDTO                     `json:"grc_report"`
	GRCReportWorkflowStatus *GRCReportWorkflowStatusDTO      `json:"grc_report_workflow_status"`
	ReviewUser              *organizationDTOs.UserMinimalDTO `json:"review_user"`
	ApproveUser             *organizationDTOs.UserMinimalDTO `json:"approve_user"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               *time.Time                       `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type GRCReportDTO struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type GRCReportWorkflowFormDTO struct {
	UserID uint   `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=risk compliance governance"`
}

type GRCReportSectionReviewApproveFormFormDTO struct {
	Type string `json:"type" binding:"required"`
}
