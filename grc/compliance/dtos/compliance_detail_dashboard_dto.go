package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceDashboardQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Page     int
	PageSize int
}

type ComplianceDetailSummaryDTO struct {
	OverallComplianceScore   float64 `json:"overall_compliance_score"`
	TotalClauses             int64   `json:"total_clauses"`
	TotalRequirements        int64   `json:"total_requirements"`
	TotalControlsImplemented int64   `json:"total_controls_implemented"`
	OverallControlEfficiency float64 `json:"overall_control_efficiency"`
}

type ComplianceClauseSummaryDTO struct {
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

type ComplianceClausePerformanceItemDTO struct {
	ComplianceClauseTitleEn string  `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr string  `json:"compliance_clause_title_ar"`
	ComplianceClauseScore   float64 `json:"compliance_clause_score"`
}

type ComplianceClausePerformanceDTO struct {
	Items []ComplianceClausePerformanceItemDTO `json:"items"`
}

type ComplianceDomainWiseOverviewItemDTO struct {
	DomainTitleEn string `json:"domain_title_en"`
	DomainTitleAr string `json:"domain_title_ar"`
	Compliant     int64  `json:"compliant"`
	Partial       int64  `json:"partial"`
	NonCompliant  int64  `json:"non_compliant"`
	Pending       int64  `json:"pending"`
}

type ComplianceDomainWiseOverviewDTO struct {
	Items      []ComplianceDomainWiseOverviewItemDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO
}
