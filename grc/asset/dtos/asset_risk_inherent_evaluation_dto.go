package dtos

type AssetRiskInherentEvaluationFormDTO struct {
	RiskMatrixLikelihoodID uint `json:"risk_matrix_likelihood_id" binding:"required"`
}

type AssetRiskMatrixLikelihoodDTO struct {
	ID            uint   `json:"id"`
	RiskMatrixID  uint   `json:"risk_matrix_id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	Score         uint   `json:"score"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Probability   string `json:"probability"`
	Status        string `json:"status"`
}

type AssetRiskMatrixImpactDTO struct {
	ID            uint   `json:"id"`
	RiskMatrixID  uint   `json:"risk_matrix_id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	Score         uint   `json:"score"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Status        string `json:"status"`
}

type AssetRiskMatrixRiskRatingDTO struct {
	ID           uint   `json:"id"`
	RiskMatrixID uint   `json:"risk_matrix_id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type AssetRiskMatrixRiskRatingValueDTO struct {
	ID                     uint `json:"id"`
	RiskMatrixRiskRatingID uint `json:"risk_matrix_risk_rating_id"`
	Value                  uint `json:"value"`
}

type AssetRiskInherentEvaluationDTO struct {
	ID                        uint                               `json:"id"`
	RiskMatrixLikelihood      *AssetRiskMatrixLikelihoodDTO      `json:"risk_matrix_likelihood"`
	RiskMatrixImpact          *AssetRiskMatrixImpactDTO          `json:"risk_matrix_impact"`
	RiskMatrixRiskRating      *AssetRiskMatrixRiskRatingDTO      `json:"risk_matrix_risk_rating"`
	RiskMatrixRiskRatingValue *AssetRiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating_value"`
	ConfidentialityScore      *int                               `json:"confidentiality_score"`
	IntegrityScore            *int                               `json:"integrity_score"`
	AvailabilityScore         *int                               `json:"availability_score"`
}
