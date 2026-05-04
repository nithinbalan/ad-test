package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	orgDTOs "jethur/grc/organization/dtos"
)

type EmailLogListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	StartDate   string
	EndDate     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type EmailLogDTO struct {
	ID                   uint                   `json:"id"`
	ModuleID             uint                   `json:"module_id"`
	SubModuleID          uint                   `json:"sub_module_id"`
	Subject              string                 `json:"subject"`
	EmailStatusID        *uint                  `json:"email_status_id"`
	EmailStatusTitleEn   *string                `json:"email_status_title_en"`
	EmailStatusTitleAr   *string                `json:"email_status_title_ar"`
	EmailStatusColorCode *string                `json:"email_status_color_code"`
	EmailStatusType      *string                `json:"email_status_type"`
	IsPermanentFailed    bool                   `json:"is_permanent_failed"`
	AttemptCount         int                    `json:"attempt_count"`
	MaxAttempts          int                    `json:"max_attempts"`
	LastAttemptAt        *time.Time             `json:"last_attempt_at"`
	ErrorMessage         *string                `json:"error_message"`
	SentAt               *time.Time             `json:"sent_at"`
	CreatedAt            time.Time              `json:"created_at"`
	Recipients           []EmailLogRecipientDTO `json:"recipients"`
	// Attachments       []EmailLogAttachmentDTO `json:"attachments"`
}

type EmailLogRecipientDTO struct {
	ID            uint                    `json:"id"`
	RecipientType string                  `json:"recipient_type"`
	EmailAddress  string                  `json:"email_address"`
	DisplayName   *string                 `json:"display_name"`
	IsExternal    bool                    `json:"is_external"`
	User          *orgDTOs.UserMinimalDTO `json:"user"`
}

type EmailLogAttachmentDTO struct {
	ID          uint      `json:"id"`
	FileName    string    `json:"file_name"`
	ContentType *string   `json:"content_type"`
	SizeBytes   *int64    `json:"size_bytes"`
	StorageKey  string    `json:"storage_key"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

type EmailLogDetailDTO struct {
	ID                  uint                       `json:"id"`
	ModuleID            uint                       `json:"module_id"`
	SubModuleID         uint                       `json:"sub_module_id"`
	Subject             string                     `json:"subject"`
	Body                string                     `json:"body"`
	EmailStatus         *masterDTOs.EmailStatusDTO `json:"email_status"`
	IsPermanentFailed   bool                       `json:"is_permanent_failed"`
	AttemptCount        int                        `json:"attempt_count"`
	MaxAttempts         int                        `json:"max_attempts"`
	LastAttemptAt       *time.Time                 `json:"last_attempt_at"`
	LastAttemptAtString *string                    `json:"last_attempt_at_string"`
	NextRetryAt         *time.Time                 `json:"next_retry_at"`
	NextRetryAtString   *string                    `json:"next_retry_at_string"`
	LockedAt            *time.Time                 `json:"locked_at"`
	LockedAtString      *string                    `json:"locked_at_string"`
	LockedBy            *string                    `json:"locked_by"`
	SmtpRemoteResponse  *string                    `json:"smtp_remote_response"`
	ErrorMessage        *string                    `json:"error_message"`
	SentAt              *time.Time                 `json:"sent_at"`
	SentAtString        *string                    `json:"sent_at_string"`
	CreatedAt           time.Time                  `json:"created_at"`
	CreatedAtString     string                     `json:"created_at_string"`
	Module              *ModuleMinimalDTO          `json:"module"`
	SubModule           *ModuleMinimalDTO          `json:"sub_module"`
	Recipients          []EmailLogRecipientDTO     `json:"recipients"`
	Attachments         []EmailLogAttachmentDTO    `json:"attachments"`
}

type EmailLogListResponseDTO struct {
	Items      []*EmailLogDTO               `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type EmailLogSummaryDTO struct {
	TotalCount      int64 `json:"total_count"`
	PendingCount    int64 `json:"pending_count"`
	SentCount       int64 `json:"sent_count"`
	ProcessingCount int64 `json:"processing_count"`
	FailedCount     int64 `json:"failed_count"`
}
