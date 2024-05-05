package dtos

type RiskMatrixRiskRatingValueFormDTO struct {
	CalculationType      string                          `json:"calculation_type" binding:"required,max=20,oneof=multiplication addition" gorm:"type:varchar(20);not null"`
	RiskMatrixRiskRating []*RiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating"`
}

type RiskMatrixRiskRatingValueDTO struct {
	ID                     uint   `json:"id"`
	RiskMatrixRiskRatingID uint   `json:"risk_matrix_risk_rating_id" binding:"required"`
	Values                 []uint `json:"values" binding:"required" gorm:"not null"`
}

type RiskMatrixRiskRatingDataDTO struct {
	RiskMatrixID           uint   `json:"risk_matrix_id"`
	RiskMatrixRiskRatingID uint   `json:"risk_matrix_risk_rating_id"`
	RiskRatingTitleEN      string `json:"risk_rating_title_en"`
	RiskRatingTitleAR      string `json:"risk_rating_title_ar"`
	RiskRatingTitleType    string `json:"risk_rating_title_type"`
	ColorCode              string `json:"color_code"`
	Value                  uint   `json:"value"`
}

type RiskMatrixRiskRatingResponseDTO struct {
	RiskMatrixImpactID          uint                        `json:"risk_matrix_impact_id"`
	RiskMatrixRiskRatingID      uint                        `json:"risk_matrix_risk_rating_id"`
	RiskMatrixImpactTitleEN     string                      `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAR     string                      `json:"risk_matrix_impact_title_ar"`
	RiskMatrixLikelihoodID      uint                        `json:"risk_matrix_likelihood_id"`
	RiskMatrixLikelihoodTitleEN string                      `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAR string                      `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixScore             uint                        `json:"risk_matrix_score"`
	RiskMatrixRiskRating        RiskMatrixRiskRatingDataDTO `json:"risk_matrix_risk_rating"`
}
