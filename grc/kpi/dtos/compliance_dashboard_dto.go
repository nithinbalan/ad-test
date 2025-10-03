package dtos

type ComplianceDashboardQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type ComplianceOverviewDTO struct {
	OrganizationComplianceScore float64 `json:"organization_compliance_score"`
	TotalComplianceRegisters    int64   `json:"total_compliance_registers"`
	TotalRequirements           int64   `json:"total_requirements"`
	TotalAssessments            int64   `json:"total_assessments"`
}

type ComplianceStatusItemDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type ComplianceBreakdownDTO struct {
	CommplainceStatuses []ComplianceStatusItemDTO `json:"compliance_statuses"`
	Total               int64                     `json:"total"`
}
