package models

import (
	"time"
)

type AuditFinalReport struct {
	ID                       uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" gorm:"type:varchar(45);not null"`
	AuditID                  uint       `json:"audit_id" gorm:"not null"`
	AuditFinalReportStatusID uint       `json:"audit_final_report_status_id" gorm:"not null"`
	SubmittedBYID            *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt              *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	NextReviewLevel          int        `json:"next_review_level" gorm:"default:0"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                    *Audit                     `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditFinalReportSections []*AuditFinalReportSection `gorm:"foreignKey:AuditFinalReportID"`
	AuditFinalReportStatus   *AuditFinalReportStatus    `gorm:"foreignKey:AuditFinalReportStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	SubmittedByUser          *User                      `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditReportWorkflows     []*AuditReportWorkflow     `gorm:"foreignKey:AuditFinalReportID"`
	CreatedByUser            *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser            *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditFinalReport to `audit_final_reports`
func (AuditFinalReport) TableName() string {
	return "audit_final_reports"
}
