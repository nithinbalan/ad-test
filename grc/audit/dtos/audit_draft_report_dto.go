package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftReportFormDTO struct {
	ReportDate       string  `json:"report_date" binding:"required"`
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}

type AuditDraftReportDetailDTO struct {
	ID                 uint                            `json:"id"`
	AuditID            uint                            `json:"audit_id"`
	Title              string                          `json:"title"`
	ReferenceNumber    string                          `json:"reference_number"`
	ReportDateString   string                          `json:"report_date"`
	PreparedBy         organizationDTOs.UserMinimalDTO `json:"prepared_by"`
	PreparedAt         string                          `json:"prepared_at"`
	ExecutiveSummary   string                          `json:"executive_summary"`
	Auditors           []*AuditPlanAuditMemberDTO      `json:"auditors"`
	TeamLead           organizationDTOs.UserMinimalDTO `json:"team_lead"`
	AuditDraftFindings []*AuditDraftFindingDetailDTO   `json:"draft_findings"`
	Conclusion         string                          `json:"conclusion"`
}
type AuditDraftReportParsedDates struct {
	ReportDate *time.Time
}
