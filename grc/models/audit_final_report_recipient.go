package models

import (
	"time"
)

// AuditFinalReportRecipient represents a recipient in the audit final report process.
type AuditFinalReportRecipient struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFinalReportID uint      `json:"audit_final_report_id" binding:"required" gorm:"not null"`
	UserID             uint      `json:"user_id" binding:"required" gorm:"not null"`
	Type               string    `json:"type" binding:"required" gorm:"type:varchar(20);default:'cc';check:type IN ('cc', 'to')"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relations
	AuditFinalReport *AuditFinalReport `gorm:"foreignKey:AuditFinalReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User             *User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the default table name for the AuditFinalReportRecipient model.
func (AuditFinalReportRecipient) TableName() string {
	return "audit_final_report_recipients"
}
