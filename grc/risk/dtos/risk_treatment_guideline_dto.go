package dtos

type RiskTreatmentGuidelineListDTO struct {
	ID                 uint   `json:"id"`
	RiskLevelID        uint   `json:"risk_level_id"`
	TimeframeEn        string `json:"timeframe_en"`
	TimeframeAr        string `json:"timeframe_ar"`
	Status             string `json:"status"`
	RiskLevelTitleEn   string `json:"risk_level_title_en"`
	RiskLevelTitleAr   string `json:"risk_level_title_ar"`
	RiskLevelColorCode string `json:"risk_level_color_code"`
}

type RiskTreatmentGuidelineFormDTO struct {
	RiskLevelID uint   `json:"risk_level_id"`
	TimeframeEn string `json:"timeframe_en"`
	TimeframeAr string `json:"timeframe_ar"`
	Status      string `json:"status"`
}

type RiskTreatmentGuidelineListQueryParamsDTO struct {
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Keyword  string              `json:"keyword"`
	SortBy   string              `json:"sort_by"`
	Sort     string              `json:"sort"`
	Status   string              `json:"status"`
	Filters  map[string][]string `json:"filters"`
}
