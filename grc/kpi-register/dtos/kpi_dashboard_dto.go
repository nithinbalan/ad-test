package dtos

import "time"

type KPIDashboardQueryParamsDTO struct {
	DomainID    *uint               `json:"domain_id"`
	Filters     map[string][]string `json:"filters"`
	Date        *time.Time          `json:"date"`
	StartDate   *time.Time          `json:"start_date"`
	EndDate     *time.Time          `json:"end_date"`
	QuickFilter []string            `json:"quick_filter"`
	PeriodType  *string             `json:"period_type"`
	Sort        string              `json:"sort"`
	SortBy      string              `json:"sort_by"`
	Keyword     string              `json:"keyword"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	IsKRI       *bool               `json:"is_kri"`
}

type KPISummaryDTO struct {
	OverallKPIAchievementPercentage float64 `json:"overall_kpi_achievement_percentage"`
	TotalActiveKPICount             int64   `json:"total_active_kpi_count"`
	PassedKPICount                  int64   `json:"passed_kpi_count"`
	FailedKPICount                  int64   `json:"failed_kpi_count"`
	OverdueKPICount                 int64   `json:"overdue_kpi_count"`
	KPIAchievedAheadOfTargetCount   int64   `json:"kpi_achieved_ahead_of_target_count"`
}

type KPIStatusCountDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type PerformanceSummaryDTO struct {
	ID                              uint    `json:"id"`
	Title                           string  `json:"title"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserRoleTitleEn      string  `json:"responsible_user_role_title_en"`
	ResponsibleUserRoleTitleAr      string  `json:"responsible_user_role_title_ar"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	PerformancePercentage           float64 `json:"performance_percentage"`
	StatusTitleEn                   string  `json:"status_title_en"`
	StatusTitleAr                   string  `json:"status_title_ar"`
	StatusColorCode                 string  `json:"status_color_code"`
	StatusType                      string  `json:"status_type"`
	PerformanceLabel                string  `json:"performance_label"`
}

type KPIPerformanceCategoryStatsDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type KPIDepartmentStatsDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type KPIPassFailSummaryResponseDTO struct {
	Total  int64 `json:"total"`
	Passed int64 `json:"passed"`
	Failed int64 `json:"failed"`
}

type KPIOverdueSummaryDTO struct {
	TotalCount         int64 `json:"total_count"`
	AverageDaysOverdue int64 `json:"average_days_overdue"`
}

type KPIOverdueReviewListDTO struct {
	ID                         uint       `json:"id"`
	Title                      string     `json:"title"`
	ResponsibleUserID          uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn *string    `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  *string    `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr *string    `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr  *string    `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial     *string    `json:"responsible_user_initial"`
	ResponsibleUserColor       *string    `json:"responsible_user_color"`
	ResponsibleUserProfilePath *string    `json:"responsible_user_profile_image_path"`
	DaysOverdue                int64      `json:"days_overdue"`
	TargetReviewDate           *time.Time `json:"target_review_date"`
}

type KPIOverdueAchievementListDTO struct {
	ID                         uint    `json:"id"`
	Title                      string  `json:"title"`
	ResponsibleUserID          uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr  *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial     *string `json:"responsible_user_initial"`
	ResponsibleUserColor       *string `json:"responsible_user_color"`
	ResponsibleUserProfilePath *string `json:"responsible_user_profile_image_path"`
	DaysOverdue                int64   `json:"days_overdue"`
	TargetValue                int64   `json:"target_value"`
	ActualValue                int64   `json:"actual_value"`
}

type KPIOverdueReviewResponseDTO struct {
	TotalCount         int64                      `json:"total_count"`
	AverageDaysOverdue int64                      `json:"average_days_overdue"`
	TopKPIs            []*KPIOverdueReviewListDTO `json:"top_kpis"`
}

type KPIOverdueAchievementResponseDTO struct {
	TotalCount         int64                           `json:"total_count"`
	AverageDaysOverdue int64                           `json:"average_days_overdue"`
	TopKPIs            []*KPIOverdueAchievementListDTO `json:"top_kpis"`
}

type KPIAchievementSummaryDTO struct {
	OverallAchievement float64 `json:"overall_achievement"`
	TotalKPICount      int64   `json:"total_kpi_count"`
}

type KPIReviewStatusDistributionDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type KPIReviewStatusDistributionResponseDTO struct {
	Total int64                             `json:"total"`
	Items []*KPIReviewStatusDistributionDTO `json:"items"`
}
type KPIAchievementByDomainDTO struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryType string  `json:"category_type"`
	AverageScore float64 `json:"average_score"`

	KPIStatusID            uint    `json:"kpi_status_id"`
	KPIStatusTitleEn       string  `json:"kpi_status_title_en"`
	KPIStatusTitleAr       string  `json:"kpi_status_title_ar"`
	KPIStatusType          string  `json:"kpi_status_type"`
	KPIStatusColorCode     string  `json:"kpi_status_color_code"`
	KPIStatusLabel         *string `json:"kpi_status_label"`
	KPIStatusOrder         *int    `json:"kpi_status_order"`
	KPIStatusIsNotAchieved bool    `json:"kpi_status_is_not_achieved"`
}

type KPIStatusOverviewDTO struct {
	ID                 uint    `json:"id"`
	TitleEn            string  `json:"title_en"`
	TitleAr            string  `json:"title_ar"`
	ColorCode          string  `json:"color_code"`
	Type               string  `json:"type"`
	Label              string  `json:"label"`
	Count              int64   `json:"count"`
	AverageAchievement float64 `json:"average_achievement"`
	AverageTrend       float64 `json:"average_trend"`
}

type KPIStatusOverviewResponseDTO struct {
	Total int64                   `json:"total"`
	Items []*KPIStatusOverviewDTO `json:"items"`
}
