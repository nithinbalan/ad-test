package dtos

type RiskInherentEvaluationFormDTO struct {
	InherentRiskMatrixLikelihoodID uint                       `json:"inherent_risk_matrix_likelihood_id"`
	InherentRiskMatrixImpactID     uint                       `json:"inherent_risk_matrix_impact_id"`
	InherentRiskMatrixRiskRatingID uint                       `json:"inherent_risk_matrix_risk_rating_id"`
	RiskImpactEvaluations          []RiskImpactEvaluationDTO  `json:"risk_impact_evaluations"`
	RiskAppetiteAlignments         []RiskAppetiteAlignmentDTO `json:"risk_appetite_alignments"`
}

type RiskImpactEvaluationDTO struct {
	RiskMatrixImpactGuidelineID uint `json:"risk_matrix_impact_guideline_id"`
	RiskImpactCategoryID        uint `json:"risk_impact_category_id"`
	RiskMatrixImpactID          uint `json:"risk_matrix_impact_id"`
}

type RiskAppetiteAlignmentDTO struct {
	RiskAppetiteStatementID      uint `json:"risk_appetite_statement_id"`
	RiskAppetiteToleranceLevelID uint `json:"risk_appetite_tolerance_level_id"`
}
