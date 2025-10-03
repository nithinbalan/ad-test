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
