package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type APIRequestLogListDTO struct {
	ID           uint       `json:"id"`
	Method       string     `json:"method"`
	URL          string     `json:"url"`
	RequestedAt  time.Time  `json:"requested_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	StatusCode   *int       `json:"status_code"`
	DurationMs   *int       `json:"duration_ms"`
	LogStatus    string     `json:"log_status"`
	RequestID    string     `json:"request_id"`
	UserID       *uint      `json:"user_id"`
	ErrorMessage *string    `json:"error_message"`
	StackTrace   *string    `json:"stack_trace"`
}

type APIRequestLogSummaryDTO struct {
	TotalRequests int64 `json:"total_requests"`
	Completed     int64 `json:"completed"`
	Pending       int64 `json:"pending"`
	Failed        int64 `json:"failed"`
	AvgDurationMs int64 `json:"avg_duration_ms"`
}

type APIRequestLogListResponseDTO struct {
	Items      []*APIRequestLogListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type APIRequestLogListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	DateFrom    string
	DateTo      string
	QuickFilter []string
}
