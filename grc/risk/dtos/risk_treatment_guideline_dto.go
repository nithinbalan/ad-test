package dtos

type RiskTreatmentGuidelineListDTO struct {
	ID                            uint   `json:"id"`
	RiskMatrixRiskRatingID        uint   `json:"risk_matrix_risk_rating_id"`
	DescriptionEn                 string `json:"description_en"`
	DescriptionAr                 string `json:"description_ar"`
	Status                        string `json:"status"`
	RiskMatrixRiskRatingTitleEn   string `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr   string `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingColorCode string `json:"risk_matrix_risk_rating_color_code"`
}

type RiskTreatmentGuidelineFormDTO struct {
	RiskMatrixRiskRatingID uint   `json:"risk_matrix_risk_rating_id"`
	DescriptionEn          string `json:"description_en"`
	DescriptionAr          string `json:"description_ar"`
	Status                 string `json:"status"`
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
