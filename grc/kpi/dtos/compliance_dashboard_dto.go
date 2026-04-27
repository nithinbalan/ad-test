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

type ComplianceLeadershipAttentionDTO struct {
	OpenCriticalRequirements      int64 `json:"open_critical_requirements"`
	OverdueComplianceRequirements int64 `json:"overdue_compliance_requirements"`
	UnassignedComplianceItems     int64 `json:"unassigned_compliance_items"`
	PastDueActions                int64 `json:"past_due_actions"`
}

type ComplianceSummaryItemDTO struct {
	ID        uint    `json:"id"`
	TitleEn   string  `json:"title_en"`
	TitleAr   string  `json:"title_ar"`
	ColorCode string  `json:"color_code"`
	Type      string  `json:"type"`
	Count     int64   `json:"count"`
	Percent   float64 `json:"percent"`
}

type ComplianceBreakdownDTO struct {
	ComplianceStatuses []ComplianceSummaryItemDTO `json:"compliance_statuses"`
	Total              int64                      `json:"total"`
}

type CompliancePerformanceRadarDTO struct {
	ComplianceStandardCode string  `json:"compliance_standard_code"`
	ComplianceTitleEn      string  `json:"compliance_title_en"`
	ComplianceTitleAr      string  `json:"compliance_title_ar"`
	Score                  float64 `json:"score"`
}

type ComplianceClausePerformanceDTO struct {
	ComplianceClauseNo      string  `json:"compliance_clause_no"`
	ComplianceClauseTitleEn string  `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr string  `json:"compliance_clause_title_ar"`
	Score                   float64 `json:"score"`
}

type CompliancePolicyOverviewDTO struct {
	TotalPolicies                  int64   `json:"total_policies"`
	AcknowledgementScore           float64 `json:"acknowledgement_score"`
	HighRiskPolicies               int64   `json:"high_risk_policies"`
	PolicyExceptionsRaised         int64   `json:"policy_exceptions_raised"`
	CriticalPoliciesUnacknowledged int64   `json:"critical_policies_unacknowledged"`
}
type CompliancePolicyStatusBreakdownDTO struct {
	TotalPolicies  int                        `json:"total_policies"`
	PolicyStatuses []ComplianceSummaryItemDTO `json:"policy_statuses"`
}

type CompliancePolicyByOrganizationUnitDTO struct {
	OrganizationUnitTitleEn   string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr   string `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode string `json:"organization_unit_color_code"`
	TotalPolicies             int64  `json:"total_policies"`
}

type ComplianceControlSummaryDTO struct {
	TotalControls          int64                      `json:"total_controls"`
	ImplementationSummary  []ComplianceSummaryItemDTO `json:"implementation_summary"`
	ImplementationStatuses []ComplianceSummaryItemDTO `json:"implementation_statuses"`
	EfficiencyChart        []ComplianceSummaryItemDTO `json:"efficiency_chart"`
	OverallEfficiency      float64                    `json:"overall_efficiency"`
}

type ComplianceEfficiencySummaryQueryResultDTO struct {
	ComplianceStandardCode string
	ComplianceTitleEn      string
	ComplianceTitleAr      string
	StatusID               uint
	TitleEn                string
	TitleAr                string
	ColorCode              string
	Type                   string
	Count                  int64
}

type ComplianceEfficiencySummaryDTO struct {
	ComplianceStandardCode string                     `json:"compliance_standard_code"`
	ComplianceTitleEn      string                     `json:"compliance_title_en"`
	ComplianceTitleAr      string                     `json:"compliance_title_ar"`
	EfficiencyStatuses     []ComplianceSummaryItemDTO `json:"efficiency_statuses"`
}

type ComplianceControlOverviewDTO struct {
	TotalAssessments       int64   `json:"total_assessments"`
	AssessmentCoverage     float64 `json:"assessment_coverage"`
	Completion             float64 `json:"completion"`
	HighRiskFindings       int64   `json:"high_risk_findings"`
	RegistersAssessedX     int64   `json:"registers_assessed_x"`
	RegistersAssessedY     int64   `json:"registers_assessed_y"`
	RegistersAssessedPct   float64 `json:"registers_assessed_pct"`
	UpcomingAssessments    int64   `json:"upcoming_assessments"`
	PendingRemediation     int64   `json:"pending_remediation"`
	DepartmentsNotAssessed int64   `json:"departments_not_assessed"`
}

type ComplianceAssessmentControlOrganizationUnit struct {
	OrganizationUnitTitleEn   string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr   string `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode string `json:"organization_unit_color_code"`
	AssessmentCount           int64  `json:"assessment_count"`
}
type ComplianceActionPlanOverviewDTO struct {
	TotalActionPlans     int64   `json:"total_action_plans"`
	CompletedPlans       int64   `json:"completed_plans"`
	PendingPlans         int64   `json:"pending_plans"`
	OverduePlans         int64   `json:"overdue_plans"`
	HighPriorityPlans    int64   `json:"high_priority_plans"`
	AvgCompletionPercent float64 `json:"avg_completion_percent"`
}

type ComplianceActionPlanStatusSummaryDTO struct {
	ActionPlanStatuses []ComplianceSummaryItemDTO `json:"action_plan_statuses"`
	TotalActionPlans   int64                      `json:"total_action_plans"`
}

type ComplianceActionPlanByOrganizationUnitDTO struct {
	OrganizationUnitTitleEn   string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr   string `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode string `json:"organization_unit_color_code"`
	TotalActionPlans          int64  `json:"total_action_plans"`
}

type ComplianceRatingDistributionDTO struct {
	Registers   []ComplianceRatingBucketDTO `json:"registers"`
	Assessments []ComplianceRatingBucketDTO `json:"assessments"`
}

type ComplianceRatingBucketDTO struct {
	Count           int64    `json:"count"`
	RatingID        uint     `json:"rating_id"`
	RatingTitleEn   string   `json:"rating_title_en"`
	RatingTitleAr   string   `json:"rating_title_ar"`
	RatingColorCode string   `json:"rating_color_code"`
	RatingType      string   `json:"rating_type"`
	RatingScoreFrom *float64 `json:"rating_score_from"`
	RatingScoreTo   *float64 `json:"rating_score_to"`
	RatingStatus    string   `json:"rating_status"`
}

type ComplianceRegisterClassificationSummaryDTO struct {
	ApplicableCount             int64 `json:"applicable_count"`
	NotApplicableCount          int64 `json:"not_applicable_count"`
	RegulationNotificationCount int64 `json:"regulation_notification_count"`
	RiskBasedAssessmentCount    int64 `json:"risk_based_assessment_count"`
	PublicCount                 int64 `json:"public_count"`
	SharedCount                 int64 `json:"shared_count"`
	RequirementsWithNoControl   int64 `json:"requirements_with_no_control"`
}
