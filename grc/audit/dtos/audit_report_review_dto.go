package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
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
type AuditFinalReportWorkflowDTO struct {
	ID                   uint                            `json:"id"`
	AuditID              uint                            `json:"audit_id"`
	AuditFinalReportID   uint                            `json:"audit_final_report_id"`
	SubmittedAt          time.Time                       `json:"submitted_at"`
	SubmittedByUser      organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	AuditReportWorkflows []*AuditReportWorkflowDTO       `json:"audit_report_workflows"`
}

type AuditReportWorkflowDTO struct {
	Level                           int                            `json:"level"`
	WorkflowStatus                  masterDTOs.WorkflowStatusDTO   `json:"workflow_status"`
	AuditReportReviews              []*AuditReportReviewDTO        `json:"audit_report_reviews"`
	AuditReportWorkflowParticipants []*WorkflowLevelParticipantDTO `json:"audit_report_workflow_participants"`
}

type AuditReportReviewDTO struct {
	Comment        string                          `json:"comment"`
	WorkflowStatus masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	CreatedByUser  organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdateByUser   organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt      time.Time                       `json:"created_at"`
	UpdatedAt      *time.Time                      `json:"updated_at"`
}

type WorkflowLevelParticipantDTO struct {
	ID   uint                             `json:"id"`
	User *organizationDTOs.UserMinimalDTO `json:"user,omitempty"`
	Role *settingsDTOs.RoleDTO            `json:"role,omitempty" `
}
