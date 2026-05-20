package models

import "time"

// APIRequestLog maps to dbo.api_request_logs (migration create_api_request_logs_table).
type APIRequestLog struct {
	ID           uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Method       string     `json:"method" gorm:"column:method;type:nvarchar(10);not null"`
	URL          string     `json:"url" gorm:"column:url;type:nvarchar(max);not null"`
	RequestID    string     `json:"request_id" gorm:"column:request_id;type:nvarchar(64);not null"`
	UserID       *uint      `json:"user_id" gorm:"column:user_id;type:int"`
	RequestBody  *string    `json:"request_body" gorm:"column:request_body;type:nvarchar(max)"`
	ResponseBody *string    `json:"response_body" gorm:"column:response_body;type:nvarchar(max)"`
	ErrorMessage *string    `json:"error_message" gorm:"column:error_message;type:nvarchar(max)"`
	StackTrace   *string    `json:"stack_trace" gorm:"column:stack_trace;type:nvarchar(max)"`
	RequestedAt  time.Time  `json:"requested_at" gorm:"column:requested_at;type:datetime;not null"`
	CompletedAt  *time.Time `json:"completed_at" gorm:"column:completed_at;type:datetime"`
	StatusCode   *int       `json:"status_code" gorm:"column:status_code"`
	DurationMs   *int       `json:"duration_ms" gorm:"column:duration_ms"`
	LogStatus    string     `json:"log_status" gorm:"column:log_status;type:nvarchar(20);not null default 'pending' check (log_status in ('pending','completed','failed','aborted'))"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (APIRequestLog) TableName() string {
	return "api_request_logs"
}
