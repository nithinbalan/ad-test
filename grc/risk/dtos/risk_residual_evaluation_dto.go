package dtos

type RiskResidualEvaluationDetailsDTO struct {
	ID                               uint                         `json:"id"`
	Type                             string                       `json:"type"`
	RiskMatrixLikelihood             RiskMatrixLikelihoodDTO      `json:"risk_matrix_likelihood"`
	RiskMatrixImpact                 RiskMatrixImpactDTO          `json:"risk_matrix_impact"`
	RiskMatrixRiskRating             RiskMatrixRiskRatingDTO      `json:"risk_matrix_risk_rating"`
	RiskMatrixRiskRatingValue        RiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating_value"`
	RiskImpactAnalyses               []RiskImpactAnalysisDTO      `json:"risk_impact_analyses"`
	RiskAppetiteAlignments           []RiskAppetiteAlignmentDTO   `json:"risk_appetite_alignments"`
	AverageControlEffectivenessScore float64                      `json:"average_risk_reduction"`
}

type RiskResidualEvaluationFormDTO struct {
	RiskMatrixLikelihoodID uint                            `json:"risk_matrix_likelihood_id"`
	RiskMatrixImpactID     uint                            `json:"risk_matrix_impact_id"`
	RiskImpactAnalyses     []*RiskImpactAnalysisFormDTO    `json:"risk_impact_analyses"`
	RiskAppetiteAlignments []*RiskAppetiteAlignmentFormDTO `json:"risk_appetite_alignments"`
}
