package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

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

type ERMDashboardListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Page     int
	PageSize int
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
	TotalRisksOverdue   int64                   `json:"total_risks_overdue"`
	OverdueByRiskLevels []ERMRiskSummaryItemDTO `json:"overdue_by_risk_levels"`
}

type ERMRiskSummaryItemDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type ERMRisksAboveAppetiteDTO struct {
	TotalRisksAboveAppetite    int64                   `json:"total_risks_above_appetite"`
	AboveAppetiteRisksByLevels []ERMRiskSummaryItemDTO `json:"above_appetite_risks_by_levels"`
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
	ControlsByEfficiency []ERMRiskSummaryItemDTO `json:"controls_by_efficiency"`
}

type ERMControlsResultDTO struct {
	TotalControls       int64
	ApplicableControls  int64
	ImplementedControls int64
	EfficientControls   int64
}

type ERMHighRiskControlsResultDTO struct {
	ControlsMitigatingHighRisks      int64
	HighRisksWithNoControls          int64
	HighRisksWithInefficientControls int64
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
	OrganizationUnit   string                  `json:"organization_unit"`
	RiskAppetiteLevels []ERMRiskSummaryItemDTO `json:"risk_appetite_levels"`
	TotalRisk          int64                   `json:"total_risk"`
}

type ERMRiskAppetiteOverviewDTO struct {
	RiskAppetiteByOrganizationUnits []ERMRiskAppetiteOverviewItemDTO `json:"risk_appetite_by_organization_units"`
}

type ERMRiskAppetiteOverviewQueryResultDTO struct {
	OrganizationUnit string `json:"organization_unit"`
	LevelID          uint   `json:"level_id"`
	TitleEn          string `json:"title_en"`
	TitleAr          string `json:"title_ar"`
	ColorCode        string `json:"color_code"`
	Type             string `json:"type"`
	Count            int64  `json:"count"`
}

type ERMRiskAppetiteSummarytDTO struct {
	RiskAppetiteLevelsSummary []ERMRiskSummaryItemDTO `json:"risk_appetite_levels_summary"`
	Total                     int64                   `json:"total"`
}

type ERMRiskAgingStatusQuereyResultDTO struct {
	Risks0To90Open         int64 `gorm:"column:risks_0_to_90_open"`
	Risks0To90Closed       int64 `gorm:"column:risks_0_to_90_closed"`
	Risks91To180Open       int64 `gorm:"column:risks_91_to_180_open"`
	Risks91To180Closed     int64 `gorm:"column:risks_91_to_180_closed"`
	Risks181To365Open      int64 `gorm:"column:risks_181_to_365_open"`
	Risks181To365Closed    int64 `gorm:"column:risks_181_to_365_closed"`
	Risks1To2YearsOpen     int64 `gorm:"column:risks_1_to_2_years_open"`
	Risks1To2YearsClosed   int64 `gorm:"column:risks_1_to_2_years_closed"`
	RisksAbove2YearsOpen   int64 `gorm:"column:risks_above_2_years_open"`
	RisksAbove2YearsClosed int64 `gorm:"column:risks_above_2_years_closed"`
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
	StrategicObjectiveTitleEn string `json:"strategic_objective_en"`
	StrategicObjectiveTitleAr string `json:"strategic_objective_ar"`
	LinkedRisks               string `json:"linked_risks"`
	HighRiskCount             int64  `json:"high_risk_count"`
	RiskRatingTitleEn         string `json:"risk_rating_title_en"`
	RiskRatingTitleAr         string `json:"risk_rating_title_ar"`
	RiskRatingColorCode       string `json:"risk_rating_color_code"`
	AppetiteBreach            string `json:"appetite_breach"`
}

type ERMStrategicObjectiveRiskOverviewDTO struct {
	Items      []ERMStregicObjectiveRiskItemDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type ERMStrategicObjectivesRiskSummaryDTO struct {
	ObjectivesWithLinkedRisks    int64 `json:"objectives_with_linked_risks"`
	ObjectivesWithHighRatedRisks int64 `json:"objectives_with_high_rated_risks"`
	ObjectivesWithNoRiskLinked   int64 `json:"objectives_with_no_risk_linked"`
}

type ERMRiskOrganizationUnitSummaryItemDTO struct {
	OrganizationUnit   string `json:"organization_unit"`
	TotalRisks         int64  `json:"total_risks"`
	HighRisks          int64  `json:"high_risks"`
	OpenRisks          int64  `json:"open_risks"`
	ClosedRisks        int64  `json:"closed_risks"`
	OverdueActionPlans int64  `json:"overdue_action_plans"`
	ResidualRiskLevel  string `json:"residual_risk_level"`
}

type ERMRiskOrganizationUnitSummaryDTO struct {
	Items      []ERMRiskOrganizationUnitSummaryItemDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
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

type ERMRiskStatusSummaryDTO struct {
	RisksByStatus []ERMRiskSummaryItemDTO `json:"risks_by_status"`
	Total         int64                   `json:"total"`
}

type ERMRiskSummaryByOrganizationUnitItemDTO struct {
	OrganizationUnit string `json:"organization_unit"`
	Count            int64  `json:"count"`
}

type ERMRiskSummaryByOrganizationUnitDTO struct {
	RisksByOrganizationUnit []ERMRiskSummaryByOrganizationUnitItemDTO `json:"risks_by_organization_unit"`
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

type ERMRiskControlEfficiencySummaryItemDTO struct {
	EfficiencyRange string  `json:"efficiency_range"`
	Percentage      float64 `json:"percentage"`
}
type ERMRiskControlEfficiencySummaryDTO struct {
	RiskControlEfficiencies []ERMRiskControlEfficiencySummaryItemDTO `json:"risk_control_efficiencies"`
}

type ERMRiskMitigationPlanStatusSummaryItemDTO struct {
	ERMRiskSummaryItemDTO
	Total int64 `json:"total"`
}

type ERMRiskMitigationPlanStatusSummaryDTO struct {
	Total                  int64                   `json:"total"`
	MitigationPlanStatuses []ERMRiskSummaryItemDTO `json:"mitigation_plan_statuses"`
}

type ERMOpenRiskSummaryItemDTO struct {
	RiskTitle                string  `json:"risk_title"`
	Category                 string  `json:"category"`
	OrganizationUnit         string  `json:"organization_unit"`
	NumControls              int64   `json:"no_of_controls"`
	ControlEfficiency        string  `json:"control_efficiency"`
	ControlEfficiencyPercent float64 `json:"control_efficiency_percent"`
	InherentRating           string  `json:"inherent_rating"`
	InherentRatingValue      int64   `json:"inherent_rating_value"`
	CurrentRating            string  `json:"current_rating"`
	CurrentRatingValue       int64   `json:"current_rating_value"`
	RiskBySource             string  `json:"risk_by_source"`
	DaysOpen                 int64   `json:"days_open"`
}

type ERMOpenRisksSummaryDTO struct {
	Items      []ERMOpenRiskSummaryItemDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
