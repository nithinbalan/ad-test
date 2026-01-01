package dtos

import "time"

type TPRMDashboardQueryParamsDTO struct {
	DateFrom string              `json:"date_from"`
	DateTo   string              `json:"date_to"`
	Filters  map[string][]string `json:"filters"`
}

type TPRMAssessmentOverviewDTO struct {
	TotalAssessmentsInitiated int64   `json:"total_assessments_initiated"`
	PublishedAssessments      int64   `json:"published_assessments"`
	AvgAssessmentScore        float64 `json:"avg_assessment_score"`
	AssessmentCompletionRate  float64 `json:"assessment_completion_rate"`
}

type TPRMThirdPartyInsightsDTO struct {
	TotalThirdParties int64 `json:"total_third_parties"`
	HighRiskCount     int64 `json:"high_risk_count"`
	OverdueCount      int64 `json:"overdue_count"`
}

type TPRMRiskLevelCountDTO struct {
	RiskLevelID uint  `json:"risk_level_id"`
	Count       int64 `json:"count"`
}

type TPRMAssessmentStatusDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type TPRMAssessmentStatusSummaryDTO struct {
	Total    int64                     `json:"total"`
	Statuses []TPRMAssessmentStatusDTO `json:"statuses"`
}

type TPRMRiskLevelDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type TPRMRiskRatingSummaryDTO struct {
	Total      int64              `json:"total"`
	RiskLevels []TPRMRiskLevelDTO `json:"risk_levels"`
}

type TPRMAssessmentScoreTrendDTO struct {
	Month string  `json:"month"`
	Score float64 `json:"score"`
}

type TPRMThirdPartyAssessmentProgressDTO struct {
	Total       int64 `json:"total"`
	Assessed    int64 `json:"assessed"`
	InProgress  int64 `json:"in_progress"`
	NotAssessed int64 `json:"not_assessed"`
}

type TPRMThirdPartyInsightsResponseDTO struct {
	AssessmentProgress TPRMThirdPartyAssessmentProgressDTO `json:"assessment_progress"`
	RiskLevels         []TPRMRiskLevelDTO                  `json:"risk_levels"`
}

type TPRMAssessmentInsightsDTO struct {
	TotalAssessments int64   `json:"total_assessments"`
	InProgress       int64   `json:"in_progress"`
	Completed        int64   `json:"completed"`
	Overdue          int64   `json:"overdue"`
	CompletionRate   float64 `json:"completion_rate"`
}

type TPRMRiskRatingDistributionDTO struct {
	Total int64              `json:"total"`
	Items []TPRMRiskLevelDTO `json:"items"`
}

type TPRMAssessmentInsightsResponseDTO struct {
	Insights               TPRMAssessmentInsightsDTO     `json:"assessment_insights"`
	RiskRatingDistribution TPRMRiskRatingDistributionDTO `json:"risk_rating_distribution"`
}

type TPRMSystemStatisticsDTO struct {
	TotalQuestions int64 `json:"total_questions"`
	Templates      int64 `json:"templates"`
}

type TPRMOngoingAssessmentDTO struct {
	ID                      string    `json:"id"`
	Title                   string    `json:"title"`
	OrganizationUnitTitleEn string    `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string    `json:"organization_unit_title_ar"`
	DueDate                 time.Time `json:"due_date"`
	DaysLeft                string    `json:"days_left"`
}

type TPRMServiceExpiryDTO struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type TPRMServiceExpiryCountsDTO struct {
	DueWithin7d  int64 `json:"due_within_7d" gorm:"column:due_within_7d"`
	OverdueLt30d int64 `json:"overdue_lt_30d" gorm:"column:overdue_lt_30d"`
	OverdueGt30d int64 `json:"overdue_gt_30d" gorm:"column:overdue_gt_30d"`
}

type TPRMAssessmentByTemplateDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type TPRMThirdPartyByTypeDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type TPRMAssessmentRiskMonthlyRowDTO struct {
	Month    string `json:"month"`
	RiskType string `json:"risk_type"`
	Count    int64  `json:"count"`
}

type TPRMAssessmentRiskMonthlyTrendDTO struct {
	Month  string `json:"month"`
	High   int64  `json:"high"`
	Medium int64  `json:"medium"`
	Low    int64  `json:"low"`
}
