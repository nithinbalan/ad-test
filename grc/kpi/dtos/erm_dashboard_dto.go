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
	RiskLevel string `json:"risk_level"`
	Count     int64  `json:"count"`
}

type ERMRisksAboveAppetiteDTO struct {
	TotalRisksAboveAppetite    int64                           `json:"total_risks_above_appetite"`
	AboveAppetiteRisksByLevels []ERMRisksAboveAppetiteLevelDTO `json:"above_appetite_risks_by_levels"`
}

type ERMRisksAboveAppetiteLevelDTO struct {
	RiskLevel string `json:"risk_level"`
	Count     int64  `json:"count"`
}

type ERMLeadershipAttentionDTO struct {
	StrategicObjectivesAtRisk int64 `json:"strategic_objectives_at_risk"`
	CriticalRisksUncontrolled int64 `json:"critical_risks_uncontrolled"`
	OverdueHighRiskActions    int64 `json:"overdue_high_risk_actions"`
	ComplianceLinkedHighRisks int64 `json:"compliance_linked_high_risks"`
}
