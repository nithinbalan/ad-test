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

type GRCIncidentSummaryDTO struct {
	IncidentClosureRate   float64 `json:"incident_closure_rate"`
	IncidentFrequencyRate float64 `json:"incident_frequency_rate"`
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

type GRCDocumentOverviewDTO struct {
	TotalDocuments int64 `json:"total_documents"`
	Policies       int64 `json:"policies"`
	Procedures     int64 `json:"procedures"`
	Others         int64 `json:"others"`
}

type GRCTopComplianceRegistersDTO struct {
	ID                        uint   `json:"id"`
	TitleEn                   string `json:"title_en"`
	TitleAr                   string `json:"title_ar"`
	TotalScore                int    `json:"total_score"`
	ComplianceStatusTitleEn   string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr   string `json:"compliance_status_title_ar"`
	ComplianceStatusType      string `json:"compliance_status_type"`
	ComplianceStatusColorCode string `json:"compliance_status_color_code"`
}

type GRCTopRiskRegistersDTO struct {
	ID                  uint   `json:"id"`
	Title               string `json:"title"`
	RiskRatingValue     uint   `json:"risk_rating_value"`
	RiskRatingTitleEn   string `json:"risk_rating_title_en"`
	RiskRatingTitleAr   string `json:"risk_rating_title_ar"`
	RiskRatingColorCode string `json:"risk_rating_color_code"`
}

type GRCTopDocumentsDTO struct {
	ID                      uint   `json:"id"`
	ReferenceNumber         string `json:"reference_number"`
	TitleEn                 string `json:"title_en"`
	TitleAr                 string `json:"title_ar"`
	DocumentType            string `json:"document_type"`
	DocumentStatusTitleEn   string `json:"document_status_title_en"`
	DocumentStatusTitleAr   string `json:"document_status_title_ar"`
	DocumentStatusColorCode string `json:"document_status_color_code"`
}

type GRCDocumentTypeChartDTO struct {
	Type       string `json:"type"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	Percentage int64  `json:"percentage"`
	Count      int64  `json:"count"`
}

type GRCDashboardListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
