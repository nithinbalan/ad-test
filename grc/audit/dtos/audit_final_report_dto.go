package dtos

type AuditFinalReportFormDTO struct {
	ExecutiveSummary *string `json:"executive_summary" binding:"required"`
	Conclusion       *string `json:"conclusion"`
}
