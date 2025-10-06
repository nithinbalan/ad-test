package dtos

type ComplianceDashboardDetailSummaryDTO struct {
	OverallComplianceScore   float64 `json:"overall_compliance_score"`
	TotalClauses             int64   `json:"total_clauses"`
	TotalRequirements        int64   `json:"total_requirements"`
	TotalControlsImplemented int64   `json:"total_controls_implemented"`
	OverallControlEfficiency float64 `json:"overall_control_efficiency"`
}

type ComplianceDashboardClauseSummaryDTO struct {
	TotalClauses                     int64   `json:"total_clauses"`
	FullyCompliantClauses            int64   `json:"fully_compliant_clauses"`
	FullyCompliantClausesPercent     float64 `json:"fully_compliant_clauses_percent"`
	PartiallyCompliantClauses        int64   `json:"partially_compliant_clauses"`
	PartiallyCompliantClausesPercent float64 `json:"partially_compliant_clauses_percent"`
	NotAssessedClauses               int64   `json:"not_assessed_clauses"`
	NotAssessedClausesPercent        float64 `json:"not_assessed_clauses_percent"`
	OverdueReviewClauses             int64   `json:"overdue_review_clauses"`
	OverdueReviewClausesPercent      float64 `json:"overdue_review_clauses_percent"`
	HighRiskClauses                  int64   `json:"high_risk_clauses"`
	HighRiskClausesPercent           float64 `json:"high_risk_clauses_percent"`
}
