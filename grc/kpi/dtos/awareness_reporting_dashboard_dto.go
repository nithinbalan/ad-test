package dtos

type AwarenessDashboardQueryParamsDTO struct {
	DateFrom string
	DateTo   string
	Year     int
}

type AwarenessPlanSummaryDTO struct {
	CampaignsPlanned       int64   `json:"campaigns_planned"`
	CampaignsWithSurvey    int64   `json:"campaigns_with_survey"`
	CampaignsWithSurveyPct float64 `json:"campaigns_with_survey_pct"`
	ActiveChannels         int64   `json:"active_channels"`
	PlannedReach           int64   `json:"planned_reach"`
}

type AwarenessCampaignsByTypeDTO struct {
	TypeTitleEn string `json:"type_title_en"`
	TypeTitleAr string `json:"type_title_ar"`
	Count       int64  `json:"count"`
}

type AwarenessPriorityDistributionDTO struct {
	PriorityTitleEn string  `json:"priority_title_en"`
	PriorityTitleAr string  `json:"priority_title_ar"`
	Count           int64   `json:"count"`
	Percentage      float64 `json:"percentage"`
}

type AwarenessCampaignsPerMonthDTO struct {
	Month int   `json:"month"`
	Count int64 `json:"count"`
}

type AwarenessChannelMixPerMonthDTO struct {
	Month           int    `json:"month"`
	ChannelTitleEn  string `json:"channel_title_en"`
	ChannelTitleAr  string `json:"channel_title_ar"`
	Count           int64  `json:"count"`
}

type AwarenessCoverageByAudienceDTO struct {
	TargetAudience string `json:"target_audience"`
	Count          int64  `json:"count"`
}

type AwarenessProgressSummaryDTO struct {
	TotalCampaigns        int64   `json:"total_campaigns"`
	CompletedCampaigns    int64   `json:"completed_campaigns"`
	OverallCompletionPct  float64 `json:"overall_completion_pct"`
	SurveyCompletionPct   float64 `json:"survey_completion_pct"`
	DeliveredOnTimePct    float64 `json:"delivered_on_time_pct"`
}


type AwarenessCompletionByChannelDTO struct {
	ChannelTitleEn string `json:"channel_title_en"`
	ChannelTitleAr string `json:"channel_title_ar"`
	Completed      int64  `json:"completed"`
	InProgress     int64  `json:"in_progress"`
	NotStarted     int64  `json:"not_started"`
}

type AwarenessCompletionPerCampaignDTO struct {
	CampaignID      uint    `json:"campaign_id"`
	ReferenceNumber string  `json:"reference_number"`
	CampaignTitle   string  `json:"campaign_title"`
	CompletionPct   float64 `json:"completion_pct"`
	StatusTitleEn   string  `json:"status_title_en"`
	StatusTitleAr   string  `json:"status_title_ar"`
}

type AwarenessSurveyCompletionDTO struct {
	TotalSurveys     int64   `json:"total_surveys"`
	CompletedSurveys int64   `json:"completed_surveys"`
	CompletionPct    float64 `json:"completion_pct"`
}

type AwarenessCompletionByDepartmentDTO struct {
	DepartmentTitleEn string  `json:"department_title_en"`
	DepartmentTitleAr string  `json:"department_title_ar"`
	CompletionPct     float64 `json:"completion_pct"`
	CampaignCount     int64   `json:"campaign_count"`
}

type AwarenessPhishingResultsDTO struct {
	AvgClickRate           float64 `json:"avg_click_rate"`
	AvgDataSubmissionRate  float64 `json:"avg_data_submission_rate"`
	AvgReportRate          float64 `json:"avg_report_rate"`
}
