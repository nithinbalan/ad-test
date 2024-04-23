package dtos

type RiskMatrixImpactDTO struct {
	ID              uint   `json:"id"`
	RiskMatrixID    uint   `json:"risk_matrix_id"`
	RiskMatrixType  string `json:"risk_matrix_type"`
	RiskMatrixTitle string `json:"risk_matrix_title"`
	TitleEn         string `json:"title_en"`
	TitleAr         string `json:"title_ar"`
	DescriptionEn   string `json:"description_en"`
	DescriptionAr   string `json:"description_ar"`
	Status          string `json:"status"`
}
