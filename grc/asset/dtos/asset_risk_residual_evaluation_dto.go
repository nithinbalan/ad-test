package dtos

type AssetRiskResidualEvaluationFormDTO struct {
	RiskMatrixLikelihoodID uint `json:"risk_matrix_likelihood_id" binding:"required"`
	RiskMatrixImpactID     uint `json:"risk_matrix_impact_id" binding:"required"`
}

type AssetRiskResidualEvaluationDTO struct {
	ID                        uint                               `json:"id"`
	RiskMatrixLikelihood      *AssetRiskMatrixLikelihoodDTO      `json:"risk_matrix_likelihood"`
	RiskMatrixImpact          *AssetRiskMatrixImpactDTO          `json:"risk_matrix_impact"`
	RiskMatrixRiskRating      *AssetRiskMatrixRiskRatingDTO      `json:"risk_matrix_risk_rating"`
	RiskMatrixRiskRatingValue *AssetRiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating_value"`
	ConfidentialityScore      *int                               `json:"confidentiality_score"`
	IntegrityScore            *int                               `json:"integrity_score"`
	AvailabilityScore         *int                               `json:"availability_score"`
}
