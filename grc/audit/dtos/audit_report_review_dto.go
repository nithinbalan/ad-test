package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditFinalReportWorkflowListDTO struct {
	ID                   uint                             `json:"id"`
	AuditID              uint                             `json:"audit_id"`
	AuditFinalReportID   uint                             `json:"audit_final_report_id"`
	SubmittedAt          *time.Time                       `json:"submitted_at"`
	SubmittedByUser      *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	AuditReportWorkflows []*AuditReportWorkflowListDTO    `json:"audit_report_workflows"`
}

type AuditReportWorkflowListDTO struct {
	Level         int                                `json:"level"`
	ApprovalType  string                             `json:"approval_type"`
	ReviewUsers   []*organizationDTOs.UserMinimalDTO `json:"review_users"`
	ApprovedUsers []*organizationDTOs.UserMinimalDTO `json:"approved_users"`
	ApprovedAt    time.Time                          `json:"approved_at,omitempty"`
	CurrentStatus masterDTOs.WorkflowStatusDTO       `json:"current_status"`
}

type AuditReportReviewCommentFormDTO struct {
	Comment    string `json:"comment"`
	IsApproved bool   `json:"is_approved"`
	IsReverted bool   `json:"is_reverted"`
	IsRejected bool   `json:"is_rejected"`
}
