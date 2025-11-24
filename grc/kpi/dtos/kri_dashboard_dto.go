package dtos

import "time"

type KriDashboardQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type KRIOverviewDTO struct {
	TotalKRIs          int64 `json:"total_kris"`
	HighRiskLinkedKRIs int64 `json:"high_risk_linked_kris"`
	HighRiskPercent    int64 `json:"high_risk_percent"`
}

type KRIReviewItemDTO struct {
	ID             uint      `json:"id"`
	NextReviewDate time.Time `json:"next_review_date"`
}

type KRIReviewStatusDTO struct {
	FourteenDaysOverdueReviews int64 `json:"fourteen_days_overdue_reviews"`
	OverdueReviews             int64 `json:"overdue_reviews"`
	UpcomingReviews            int64 `json:"upcoming_reviews"`
}

type KriBreachRawDTO struct {
	KriID      uint
	KriTitle   string
	Value      int
	ReviewedAt time.Time
	LowTh      *int
	MedTh      *int
}

type KriBreachItemDTO struct {
	Exceeded int    `json:"exceeded"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}

type KriMostFrequentBreachedResponseDTO struct {
	KriID    uint               `json:"kri_id"`
	KriTitle string             `json:"kri_title"`
	Breaches []KriBreachItemDTO `json:"breaches"`
}

type KRIStatusCategoryItemDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type KRIStatusCategoryDTO struct {
	Total      int64                      `json:"total"`
	Categories []KRIStatusCategoryItemDTO `json:"categories"`
}

type KRISeverityItemDTO struct {
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type KRISeverityBreakdownDTO struct {
	Total      int64                `json:"total"`
	Severities []KRISeverityItemDTO `json:"severities"`
}

type KriExposureRaw struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Score           *int   `json:"score"`
	LowThreshold    *int   `json:"low_threshold"`
	MediumThreshold *int   `json:"medium_threshold"`
}

type KriRiskExposureDTO struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Score     int     `json:"score"`
	Threshold int     `json:"threshold"`
	Exposure  float64 `json:"exposure"`
}

type KriOrganizationUnitsRawDTO struct {
	OrganizationUnitID      uint
	OrganizationUnitTitleEn string
	OrganizationUnitTitleAr string
	RiskLevelTitleEn        string
	RiskLevelTitleAr        string
	RiskLevelColorCode      string
	RiskLevelType           string
	Count                   int64
}

type KriSeverityOrganizationUnitDTO struct {
	OrganizationUnitTitleEn string               `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string               `json:"organization_unit_title_ar"`
	Severities              []KRISeverityItemDTO `json:"severities"`
}

type KriCategorySeverityRawDTO struct {
	CategoryID         uint
	CategoryTitleEn    string
	CategoryTitleAr    string
	RiskLevelTitleEn   string
	RiskLevelTitleAr   string
	RiskLevelColorCode string
	RiskLevelType      string
	Count              int64
}

type KriSeverityCategoryDTO struct {
	CategoryTitleEn string               `json:"category_title_en"`
	CategoryTitleAr string               `json:"category_title_ar"`
	Severities      []KRISeverityItemDTO `json:"severities"`
}
