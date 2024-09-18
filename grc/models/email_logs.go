package models

import "time"

type EmailLog struct {
	ID             uint   `gorm:"primaryKey"`
	RecipientEmail string `gorm:"type:varchar(255);not null"`
	Subject        string `gorm:"type:varchar(255);not null"`
	TemplatePath   string `gorm:"type:varchar(255);not null"`
	EmailData      string `gorm:"type:text"`
	SendStatus     string `gorm:"type:varchar(50);not null"`
	ErrorMessage   string `gorm:"type:text"`
	SentAt         *time.Time
	CreatedByID    uint      `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

// TableName overrides the table name used by ActionPlanStatus to `action_plan_statuses`
func (EmailLog) TableName() string {
	return "email_logs"
}
