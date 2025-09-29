package dtos

type ERMRiskOverviewDTO struct {
	TotalRisks          int64 `json:"total_risk"`
	OpenRisks           int64 `json:"open_risks"`
	ClosedRisks         int64 `json:"closed_risks"`
	UnresolvedHighRisks int64 `json:"unresolved_high_risks"`
}

type ERMDashboardQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type ERMRiskDistributionDTO struct {
	RiskCountsByRating []ERMRiskRatingCountDTO `json:"risk_counts_by_rating"`
}

type ERMRiskRatingCountDTO struct {
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	OpenCount  int64   `json:"open_count"`
	Percentage float64 `json:"percentage"`
}

type ERMOverdueRisksDTO struct {
	TotalRisksOverdue   int64                     `json:"total_risks_overdue"`
	OverdueByRiskLevels []ERMOverdueRiskDetailDTO `json:"overdue_by_risk_levels"`
}

type ERMOverdueRiskDetailDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type ERMRisksAboveAppetiteDTO struct {
	TotalRisksAboveAppetite    int64                           `json:"total_risks_above_appetite"`
	AboveAppetiteRisksByLevels []ERMRisksAboveAppetiteLevelDTO `json:"above_appetite_risks_by_levels"`
}

type ERMRisksAboveAppetiteLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type ERMLeadershipAttentionDTO struct {
	StrategicObjectivesAtRisk int64 `json:"strategic_objectives_at_risk"`
	CriticalRisksUncontrolled int64 `json:"critical_risks_uncontrolled"`
	OverdueHighRiskActions    int64 `json:"overdue_high_risk_actions"`
	ComplianceLinkedHighRisks int64 `json:"compliance_linked_high_risks"`
}

type ERMControlEfficiencySummaryDTO struct {
	OverallControlEfficiency float64 `json:"overall_control_efficiency"`
}

type ERMControlEfficiencyBreakdownDTO struct {
	ControlsByEfficiency []ERMControlEfficiencyByStatusDTO `json:"controls_by_efficiency"`
}

type ERMControlEfficiencyByStatusDTO struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int64  `json:"count"`
}

type ERMControlsOverviewDTO struct {
	TotalControls                    int64   `json:"total_controls"`
	TotalControlsPercent             float64 `json:"total_controls_percent"`
	ApplicableControls               int64   `json:"applicable_controls"`
	ApplicableControlsPercent        float64 `json:"applicable_controls_percent"`
	ImplementedControls              int64   `json:"implemented_controls"`
	ImplementedControlsPercent       float64 `json:"implemented_controls_percent"`
	OverallControlEfficiency         float64 `json:"overall_control_efficiency"`
	ControlsMitigatingHighRisks      int64   `json:"controls_mitigating_high_risks"`
	HighRisksWithNoControls          int64   `json:"high_risks_with_no_controls"`
	HighRisksWithInefficientControls int64   `json:"high_risks_with_inefficient_controls"`
	HighRisksWithInefficientPercent  float64 `json:"high_risks_with_inefficient_percent"`
}

type ERMRiskAppetiteOverviewItemDTO struct {
	Department string `json:"department"`
	Above      int64  `json:"above"`
	Within     int64  `json:"within"`
	Below      int64  `json:"below"`
	TotalRisk  int64  `json:"total_risk"`
}

type ERMRiskAppetiteOverviewDTO struct {
	RiskAppetiteByDepartments []ERMRiskAppetiteOverviewItemDTO `json:"risk_appetite_by_departments"`
}

type ERMRiskAppetiteLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type ERMRiskAppetiteSummarytDTO struct {
	RiskAppetiteLevelsSummary []ERMRiskAppetiteLevelDTO `json:"risk_appetite_levels_summary"`
	Total                     int64                     `json:"total"`
}

type ERMRiskAgingStatusDTO struct {
	RiskAgingPeriods []RiskAgingPeriodDTO `json:"risk_aging_periods"`
}

type RiskAgingPeriodDTO struct {
	RiskAging string `json:"risk_aging"`
	Open      int64  `json:"open"`
	Closed    int64  `json:"closed"`
	Total     int64  `json:"total"`
}

type ERMRiskTreatmentPlanProgressItemDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	OnTime    int64  `json:"on_time"`
	Overdue   int64  `json:"overdue"`
}

type ERMRiskTreatmentPlanProgressDTO struct {
	RiskTreatmentPlanProgress []ERMRiskTreatmentPlanProgressItemDTO `json:"risk_treatment_plan_progress"`
}

type ERMStregicObjectiveRiskItemDTO struct {
	StrategicObjective string `json:"strategic_objective"`
	LinkedRisks        string `json:"linked_risks"`
	HighRiskCount      int64  `json:"high_risk_count"`
	RiskRating         string `json:"risk_rating"`
	AppetiteBreach     string `json:"appetite_breach"`
}

type ERMStrategicObjectiveRiskOverviewDTO struct {
	StrategicObjectiveRisks []ERMStregicObjectiveRiskItemDTO `json:"strategic_objective_risks"`
}

type ERMStrategicObjectivesRiskSummaryDTO struct {
	ObjectivesWithLinkedRisks    int64 `json:"objectives_with_linked_risks"`
	ObjectivesWithHighRatedRisks int64 `json:"objectives_with_high_rated_risks"`
	ObjectivesWithNoRiskLinked   int64 `json:"objectives_with_no_risk_linked"`
}

type ERMRiskDepartmentSummaryItemDTO struct {
	Department         string `json:"department"`
	TotalRisks         int64  `json:"total_risks"`
	HighRisks          int64  `json:"high_risks"`
	OpenRisks          int64  `json:"open_risks"`
	ClosedRisks        int64  `json:"closed_risks"`
	OverdueActionPlans int64  `json:"overdue_action_plans"`
	ResidualRiskLevel  string `json:"residual_risk_level"`
}

type ERMRiskDepartmentSummaryDTO struct {
	RiskSummaryByDepartments []ERMRiskDepartmentSummaryItemDTO `json:"risk_summary_by_departments"`
}

// ERM Risk Details
type ERMRiskCategorySummaryItemDTO struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type ERMRiskCategorySummaryDTO struct {
	RisksByCategory []ERMRiskCategorySummaryItemDTO `json:"risks_by_category"`
}

type ERMRiskSourceSummaryItemDTO struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

type ERMRiskSourceSummaryDTO struct {
	RisksBySource []ERMRiskSourceSummaryItemDTO `json:"risks_by_source"`
}

type ERMRiskStatusSummaryItemDTO struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type ERMRiskStatusSummaryDTO struct {
	RisksByStatus []ERMRiskStatusSummaryItemDTO `json:"risks_by_status"`
	Total         int64                         `json:"total"`
}

type ERMRiskSummaryByDepartmentItemDTO struct {
	Department string `json:"department"`
	Count      int64  `json:"count"`
}

type ERMRiskSummaryByDepartmentDTO struct {
	RisksByDepartment []ERMRiskSummaryByDepartmentItemDTO `json:"risks_by_department"`
}

type ERMRisksIdentifiedSummaryItemDTO struct {
	Year  int64 `json:"year"`
	Count int64 `json:"count"`
}

type ERMRisksIdentifiedSummaryDTO struct {
	RisksIdentifiedOverYears []ERMRisksIdentifiedSummaryItemDTO `json:"risks_identified_over_years"`
}

type ERMRiskResponseStrategySummaryItemDTO struct {
	ResponseStrategy string `json:"response_strategy"`
	Count            int64  `json:"count"`
}

type ERMRiskResponseStrategySummaryDTO struct {
	RisksResponseStrategySummary []ERMRiskResponseStrategySummaryItemDTO `json:"risks_by_response_strategy"`
}

type RiskControlEfficiencySummaryItemDTO struct {
	EfficiencyRange string  `json:"efficiency_range"`
	Percentage      float64 `json:"percentage"`
}
type RiskControlEfficiencySummaryDTO struct {
	RiskControlEfficiency []RiskControlEfficiencySummaryItemDTO `json:"risk_control_efficiency"`
}
