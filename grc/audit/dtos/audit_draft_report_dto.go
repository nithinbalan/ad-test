package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftReportFormDTO struct {
	ReportDate       string  `json:"report_date" binding:"required"`
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditDraftReportDetailDTO struct {
	ID                 uint                                `json:"id"`
	AuditID            uint                                `json:"audit_id"`
	Title              string                              `json:"title"`
	ReferenceNumber    string                              `json:"reference_number"`
	ReportDateString   string                              `json:"report_date"`
	PreparedAt         string                              `json:"prepared_at"`
	PreparedBy         organizationDTOs.UserMinimalDTO     `json:"prepared_by"`
	ExecutiveSummary   string                              `json:"executive_summary"`
	Auditors           []*AuditPlanAuditMemberDTO          `json:"auditors"`
	TeamLead           organizationDTOs.UserMinimalDTO     `json:"team_lead"`
	AuditDraftFindings []*AuditDraftFindingReportDetailDTO `json:"draft_findings"`
	Conclusion         string                              `json:"conclusion"`
}
type AuditDraftReportParsedDates struct {
	ReportDate *time.Time
}
type AuditDraftFindingReportDetailDTO struct {
	ID                        uint                             `json:"id"`
	Title                     string                           `json:"title"`
	Description               *string                          `json:"description"`
	AuditID                   uint                             `json:"audit_id"`
	AuditableItemID           *uint                            `json:"auditable_item_id"`
	AuditChecklistID          *uint                            `json:"audit_checklist_id"`
	AuditFindingCategoryID    uint                             `json:"audit_finding_category_id"`
	SeverityLevelID           uint                             `json:"severity_level_id"`
	Recommendations           *string                          `json:"recommendations,omitempty"`
	AuditDraftFindingStatusID uint                             `json:"audit_draft_finding_status_id"`
	AuditableItem             AuditableItemDetailDTO           `json:"auditable_item"`
	AuditFindingCategory      AuditFindingCategoryDTO          `json:"audit_finding_category"`
	SeverityLevel             masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	AuditDraftFindingStatus   AuditDraftFindingStatusDTO       `json:"audit_draft_finding_status"`
	CreatedAt                 time.Time                        `json:"created_at"`
	UpdatedAt                 *time.Time                       `json:"updated_at"`
	CreatedByUser             organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser             *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Documents                 []*AuditDraftFindingDocumentDTO  `json:"documents"`
}
