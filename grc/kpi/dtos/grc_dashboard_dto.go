package dtos

type GRCComplianceSummaryDTO struct {
	ComplianceRequirementFulfillmentRate float64 `json:"compliance_requirement_fulfillment_rate"`
	ComplianceActionPlanResolutionRate   float64 `json:"compliance_action_plan_resolution_rate"`
}

type GRCRiskSummaryDTO struct {
	CountOfOpenHighRiskRatings int64   `json:"count_of_open_high_risk_ratings"`
	RiskClosureRate            float64 `json:"risk_closure_rate"`
}

type GRCAuditSummaryDTO struct {
	AuditCompletionRate        float64 `json:"audit_completion_rate"`
	OpenMajorNCCount           int64   `json:"open_major_nc_count"`
	AuditFindingResolutionRate float64 `json:"audit_finding_resolution_rate"`
}

type GRCDocumentSummaryDTO struct {
	TotalPolicies       int64 `json:"total_policies"`
	TotalOtherDocuments int64 `json:"total_other_documents"`
}

type GRCOrganizationSummaryDTO struct {
	ControlEffectivenessRate float64 `json:"control_effectiveness_rate"`
}

type GRCComplianceOverviewDTO struct {
	TotalComplaintsRegistered int64 `json:"total_complaints_registered"`
	Compliant                 int64 `json:"compliant"`
	NotCompliant              int64 `json:"not_compliant"`
	PartiallyCompliant        int64 `json:"partially_compliant"`
	Others                    int64 `json:"others"`
}

type GRCRiskOverviewDTO struct {
	TotalRisksRegistered int64            `json:"total_risks_registered"`
	RiskCountsByRating   map[string]int64 `json:"risk_counts_by_rating"`
}

type GRCDashboardListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
