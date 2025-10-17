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
	TotalControls                           int64   `json:"total_controls"`
	EfficientControls                       int64   `json:"efficient_controls"`
	EfficientControlsPercent                float64 `json:"efficient_controls_percent"`
	ImplementedControls                     int64   `json:"implemented_controls"`
	ImplementedControlsPercent              float64 `json:"implemented_controls_percent"`
	OverallControlEfficiency                float64 `json:"overall_control_efficiency"`
	ControlsMitigatingHighRisks             int64   `json:"controls_mitigating_high_risks"`
	HighRisksWithNoControls                 int64   `json:"high_risks_with_no_controls"`
	HighRisksWithInefficientControls        int64   `json:"high_risks_with_inefficient_controls"`
	HighRisksWithInefficientControlsPercent float64 `json:"high_risks_with_inefficient_controls_percent"`
}

type ERMRiskAppetiteOverviewDTO struct {
	OrganizationUnit   string                  `json:"organization_unit"`
	RiskAppetiteLevels []ERMRiskSummaryItemDTO `json:"risk_appetite_levels"`
	TotalRisk          int64                   `json:"total_risk"`
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
	RiskAging string `json:"risk_aging"`
	Open      int64  `json:"open"`
	Closed    int64  `json:"closed"`
	Total     int64  `json:"total"`
}

type ERMRiskTreatmentPlanProgressDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	OnTime    int64  `json:"on_time"`
	Overdue   int64  `json:"overdue"`
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
	OrganizationUnitTitleEn     string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr     string `json:"organization_unit_title_ar"`
	TotalRisks                  int64  `json:"total_risks"`
	HighRisks                   int64  `json:"high_risks"`
	OpenRisks                   int64  `json:"open_risks"`
	ClosedRisks                 int64  `json:"closed_risks"`
	OverdueActionPlans          int64  `json:"overdue_action_plans"`
	ResidualRiskRatingTitleEn   string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr   string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode string `json:"residual_risk_rating_color_code"`
}

type ERMRiskOrganizationUnitSummaryDTO struct {
	Items      []ERMRiskOrganizationUnitSummaryItemDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

// ERM Risk Details
type ERMRiskCategorySummaryItemDTO struct {
	CategoryTitleEn string `json:"category_title_en"`
	CategoryTitleAr string `json:"category_title_ar"`
	Count           int64  `json:"count"`
}

type ERMRiskCategorySummaryDTO struct {
	RisksByCategory []ERMRiskCategorySummaryItemDTO `json:"risks_by_category"`
	Total           int64                           `json:"total"`
}

type ERMRiskSourceSummaryDTO struct {
	SourceTitleEn string `json:"source_title_en"`
	SourceTitleAR string `json:"source_title_ar"`
	Count         int64  `json:"count"`
}
type ERMRiskStatusSummaryDTO struct {
	RisksByStatus []ERMRiskSummaryItemDTO `json:"risks_by_status"`
	Total         int64                   `json:"total"`
}

type ERMRiskSummaryByOrganizationUnitDTO struct {
	OrganizationUnitTitleEn   string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr   string `json:"organization_unit_title_ar"`
	OrganizationUnitColorCode string `json:"organization_unit_color_code"`
	Count                     int64  `json:"count"`
}
type ERMRisksIdentifiedSummaryDTO struct {
	Year  int64 `json:"year"`
	Count int64 `json:"count"`
}

type ERMRiskResponseStrategySummaryDTO struct {
	ResponseStrategyTitleEn string `json:"response_strategy_title_en"`
	ResponseStrategyTitleAr string `json:"response_strategy_title_ar"`
	Count                   int64  `json:"count"`
}

type ERMRiskControlEfficiencySummaryDTO struct {
	EfficiencyRange string  `json:"efficiency_range"`
	Percentage      float64 `json:"percentage"`
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
	RiskTitle                  string  `json:"risk_title"`
	Category                   string  `json:"category"`
	OrganizationUnit           string  `json:"organization_unit"`
	NumControls                int64   `json:"no_of_controls"`
	ControlEfficiency          string  `json:"control_efficiency"`
	ControlEfficiencyPercent   float64 `json:"control_efficiency_percent"`
	InherentRating             string  `json:"inherent_rating"`
	InherentRatingValue        int64   `json:"inherent_rating_value"`
	CurrentRating              string  `json:"current_rating"`
	CurrentRatingValue         int64   `json:"current_rating_value"`
	RiskAppetiteLevelTitleEn   string  `json:"risk_appetite_level_title_en"`
	RiskAppetiteLevelTitleAr   string  `json:"risk_appetite_level_title_ar"`
	RiskAppetiteLevelColorCode string  `json:"risk_appetite_level_color_code"`
	DaysOpen                   int64   `json:"days_open"`
}

type ERMOpenRisksSummaryDTO struct {
	Items      []ERMOpenRiskSummaryItemDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
