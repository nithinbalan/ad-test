package models

import (
	"time"
)

type LogSettings struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditLogRetentionDays  int        `json:"audit_log_retention_days" binding:"required" gorm:"not null;default:0"`
	EmailLogRetentionDays  int        `json:"email_log_retention_days" binding:"required" gorm:"not null;default:0"`
	SystemLogRetentionDays int        `json:"system_log_retention_days" binding:"required" gorm:"not null;default:0"`
	UpdatedByID            *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by LogSettings to `log_settings`
func (LogSettings) TableName() string {
	return "log_settings"
}
