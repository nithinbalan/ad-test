package models

import (
	"time"
)

// ErrorLog represents the error_logs table in GORM
type ErrorLog struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	HTTPStatusCode    int       `json:"http_status_code" validate:"required,numeric"`
	HTTPRequestMethod string    `gorm:"size:10" json:"http_request_method" validate:"required,oneof=GET POST PUT DELETE PATCH OPTIONS"`
	URL               string    `gorm:"size:2048" json:"url" validate:"required,url"`
	ErrorMessage      string    `gorm:"type:text" json:"error_message" validate:"required"`
	StackTrace        string    `gorm:"type:text" json:"stack_trace"`
	UserID            *uint     `json:"user_id" validate:"omitempty,numeric"`
	CreatedAt         time.Time `json:"created_at" validate:"required"`
}

// TableName overrides the table name
func (ErrorLog) TableName() string {
	return "error_logs"
}
