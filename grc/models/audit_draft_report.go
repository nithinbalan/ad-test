package models

import (
	"time"
)

type AuditDraftReport struct {
	ID                       uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID                  uint       `json:"audit_id" gorm:"not null"`
	ReportDate               time.Time  `json:"report_date" gorm:"type:date;not null"`
	ExecutiveSummary         *string    `json:"executive_summary" gorm:"type:nvarchar(MAX);default:null"`
	Conclusion               *string    `json:"conclusion" gorm:"type:nvarchar(MAX);default:null"`
	SubmittedBYID            *uint      `json:"submitted_by_id" gorm:"default:null"`                    
	SubmittedAt              *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`         
	AuditDraftReportStatusID uint       `json:"audit_draft_report_status_id" gorm:"not null;"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                  *Audit                  `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	SubmittedByUser        *User                   `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`            
	AuditDraftReportStatus *AuditDraftReportStatus `gorm:"foreignKey:AuditDraftReportStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"` 
}

// TableName overrides the table name used by AuditDraftReport to `audit_draft_reports`
func (AuditDraftReport) TableName() string {
	return "audit_draft_reports"
}
