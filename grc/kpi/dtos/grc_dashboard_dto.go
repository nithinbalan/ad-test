package dtos

type GRCComplianceSummaryDTO struct {
	ComplianceRequirementFulfillmentRate float64 `json:"compliance_requirement_fulfillment_rate"`
	ComplianceActionPlanResolutionRate   float64 `json:"compliance_action_plan_resolution_rate"`
}

type GRCRiskSummaryDTO struct {
	CountOfOpenHighRiskRatings int64   `json:"count_of_open_high_risk_ratings"`
	RiskClosureRate            float64 `json:"risk_closure_rate"`
}

type GRCDashboardListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
