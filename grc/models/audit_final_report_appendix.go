package models

import (
	"time"
)

type AuditFinalReportAppendix struct {
	ID                 uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditFinalReportID uint       `json:"audit_final_report_id" gorm:"not null"`
	URL                string     `json:"url"`
	Name               string     `json:"name" binding:"required"`
	Size               int        `json:"size"`
	Token              string     `json:"token"`
	OriginalName       string     `json:"original_name"`
	Type               string     `json:"type"`
	CreatedBYID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AuditFinalReport *AuditFinalReport `gorm:"foreignKey:AuditFinalReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditFinalReportAppendix
func (AuditFinalReportAppendix) TableName() string {
	return "audit_final_report_appendices"
}
