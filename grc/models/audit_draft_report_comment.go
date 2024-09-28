package models

import (
	"time"
)

type AuditDraftReportComment struct {
	ID                        uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditDraftReportID        uint       `json:"audit_draft_report_id" gorm:"not null"`
	AuditDraftFindingID       *uint      `json:"audit_draft_finding_id" gorm:"default:null"`
	AuditDraftReportCommentID *uint      `json:"audit_draft_report_comment_id" gorm:"default:null"`
	Type                      string     `json:"type" gorm:"type:varchar(50);not null;check:type IN ('executive-summary', 'finding', 'conclusion')"`
	Comment                   string     `json:"comment" binding:"required" gorm:"type:nvarchar(MAX);not null"`
	IsResolved                bool       `json:"is_resolved" gorm:"type:bit;default:0;not null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID               uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditDraftReport  *AuditDraftReport                  `gorm:"foreignKey:AuditDraftReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditDraftFinding *AuditDraftFinding                 `gorm:"foreignKey:AuditDraftFindingID"`
	ParentComment     *AuditDraftReportComment           `gorm:"foreignKey:AuditDraftReportCommentID"`
	ChildComments     []*AuditDraftReportComment         `gorm:"foreignKey:AuditDraftReportCommentID"`
	Documents         []*AuditDraftReportCommentDocument `gorm:"foreignKey:AuditDraftReportCommentID"`
	CreatedByUser     *User                              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser     *User                              `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditDraftReportComment to `audit_draft_report_comments`
func (AuditDraftReportComment) TableName() string {
	return "audit_draft_report_comments"
}
