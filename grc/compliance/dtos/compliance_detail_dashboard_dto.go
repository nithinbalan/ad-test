package dtos

type ComplianceDashboardDetailSummaryDTO struct {
	OverallComplianceScore   float64 `json:"overall_compliance_score"`
	TotalClauses             int64   `json:"total_clauses"`
	TotalRequirements        int64   `json:"total_requirements"`
	TotalControlsImplemented int64   `json:"total_controls_implemented"`
	OverallControlEfficiency float64 `json:"overall_control_efficiency"`
}
