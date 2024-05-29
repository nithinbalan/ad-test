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

type CreateTaskDTO struct {
	ID                     *uint      `json:"id,omitempty"`
	ModuleID               uint       `json:"module_id" binding:"required"`
	ComplianceActionPlanID *uint      `json:"compliance_action_plan_id,omitempty"`
	SLAActionPlanID        *uint      `json:"sla_action_plan_id,omitempty"`
	IncidentActionPlanID   *uint      `json:"incident_action_plan_id,omitempty"`
	ControlActionPlanID    *uint      `json:"control_action_plan_id,omitempty"`
	Title                  string     `json:"title" binding:"required"`
	Description            *string    `json:"description,omitempty"`
	TargetDate             *time.Time `json:"target_date,omitempty"`
	URL                    string     `json:"url" binding:"required"`
	SeverityLevelID        *uint      `json:"severity_level_id,omitempty"`
	ResponsibleUserID      *uint      `json:"responsible_user_id,omitempty"`
	TaskStatusID           *uint      `json:"task_status_id"`
}
