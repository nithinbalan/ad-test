package dtos

import "time"

type PaginationInfoDTO struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type SuccessfulIDsResponse struct {
	SuccessfulIds []uint `json:"successfull_ids"`
}

type FailedIDsResponse struct {
	FailedIds []uint `json:"failed_ids"`
}

type FileDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FileType     string `json:"type"`
	SizeInBytes  int64  `json:"size"`
}

type FilePreviewDTO struct {
	URL   string `json:"url" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type DocumentFormDTO struct {
	ID            *uint   `json:"id,omitempty"`
	URL           *string `json:"url"`
	FileName      string  `json:"file_name"`
	OriginalName  string  `json:"original_name"`
	Type          string  `json:"type"`
	Size          int     `json:"size"`
	IsNewFile     *bool   `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool   `json:"is_deleted_file,omitempty"`
}

type TaskFormDTO struct {
	ID                       *uint      `json:"id,omitempty"`
	ModuleIdentifier         string     `json:"module_identifier" binding:"required"`
	ComplianceActionPlanID   *uint      `json:"compliance_action_plan_id,omitempty"`
	SLAActionPlanID          *uint      `json:"sla_action_plan_id,omitempty"`
	IncidentActionPlanID     *uint      `json:"incident_action_plan_id,omitempty"`
	ControlActionPlanID      *uint      `json:"control_action_plan_id,omitempty"`
	RiskActionPlanID         *uint      `json:"risk_action_plan_id,omitempty"`
	MeetingActionPlanID      *uint      `json:"meeting_action_plan_id,omitempty"`
	AuditFindingActionPlanID *uint      `json:"audit_finding_action_plan_id,omitempty"`
	DocumentID               *uint      `json:"document_id,omitempty"`
	DocumentWorkflowID       *uint      `json:"document_workflow_id,omitempty"`
	AuditFinalReportID       *uint      `json:"audit_final_report_id,omitempty"`
	AuditDraftReportID       *uint      `json:"audit_draft_report_id,omitempty"`
	AuditReportWorkflowID    *uint      `json:"audit_report_workflow_id,omitempty"`
	Title                    string     `json:"title" binding:"required"`
	Description              *string    `json:"description,omitempty"`
	TargetDate               *time.Time `json:"target_date,omitempty"`
	URLParam                 *string    `json:"url_param,omitempty"`
	QueryParam               *string    `json:"query_param,omitempty"`
	SeverityLevelID          *uint      `json:"severity_level_id,omitempty"`
	ResponsibleUserID        uint       `json:"responsible_user_id" binding:"required"`
	TaskStatusID             *uint      `json:"task_status_id"`
}
