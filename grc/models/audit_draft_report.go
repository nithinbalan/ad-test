package models

import (
	"time"
)

type AuditDraftReport struct {
	ID                          uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID                     uint       `json:"audit_id" gorm:"not null"`
	ReportDate                  time.Time  `json:"report_date" gorm:"type:date;not null"`
	ExecutiveSummary            *string    `json:"executive_summary" gorm:"type:nvarchar(MAX);default:null"`
	Conclusion                  *string    `json:"conclusion" gorm:"type:nvarchar(MAX);default:null"`
	InternalReviewSubmittedBYID *uint      `json:"internal_review_submitted_by_id" gorm:"default:null"`
	InternalReviewSubmittedAt   *time.Time `json:"internal_review_submitted_at" gorm:"type:datetime;default:null"`
	AuditeeReviewSubmittedBYID  *uint      `json:"auditee_review_submitted_by_id" gorm:"default:null"`
	AuditeeReviewSubmittedAt    *time.Time `json:"auditee_review_submitted_at" gorm:"type:datetime;default:null"`
	AuditDraftReportStatusID    uint       `json:"audit_draft_report_status_id" gorm:"not null;"`
	NextReviewLevel             int        `json:"next_review_level" gorm:"default:0"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                         *Audit                     `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser                 *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                 *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	InternalReviewSubmittedByUser *User                      `gorm:"foreignKey:InternalReviewSubmittedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditeeReviewSubmittedByUser  *User                      `gorm:"foreignKey:AuditeeReviewSubmittedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditDraftReportStatus        *AuditDraftReportStatus    `gorm:"foreignKey:AuditDraftReportStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditDraftReportComments      []*AuditDraftReportComment `gorm:"foreignKey:AuditDraftReportID"`
}

// TableName overrides the table name used by AuditDraftReport to `audit_draft_reports`
func (AuditDraftReport) TableName() string {
	return "audit_draft_reports"
}
