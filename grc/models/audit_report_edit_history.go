package models

import (
	"time"
)

type AuditReportEditHistory struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement" binding:"required"`
	AuditID            uint      `json:"document_id" gorm:"not null;index" binding:"required"`
	AuditDraftReportID *uint     `json:"audit_draft_report_id" gorm:"default:null"`
	AuditFinalReportID *uint     `json:"audit_final_report_id" gorm:"default:null"`
	ActionType         string    `json:"action_type" gorm:"type:varchar(20);not null;default:'create'" binding:"required,oneof=create update submit revert approve publish"`
	Comment            *string   `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID        uint      `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null" binding:"required"`

	// Relations
	Audit            *Audit            `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AuditDraftReport *AuditDraftReport `gorm:"foreignKey:AuditDraftReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditFinalReport *AuditFinalReport `gorm:"foreignKey:AuditFinalReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentEditHistory) TableName() string {
	return "document_edit_histories"
}
