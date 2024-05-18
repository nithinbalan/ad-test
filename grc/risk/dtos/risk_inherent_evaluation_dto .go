package dtos

type RiskInherentEvaluationDetailsDTO struct {
	ID                           int                         `json:"id"`
	InherentRiskMatrixLikelihood RiskMatrixLikelihoodDTO     `json:"inherent_risk_matrix_likelihood"`
	InherentRiskMatrixImpact     RiskMatrixImpactDTO         `json:"inherent_risk_matrix_impact"`
	InherentRiskMatrixRiskRating RiskMatrixRiskRatingDTO     `json:"inherent_risk_matrix_risk_rating"`
	RiskImpactEvaluations        []*RiskImpactEvaluationDTO  `json:"risk_impact_evaluations"`
	RiskAppetiteAlignments       []*RiskAppetiteAlignmentDTO `json:"risk_appetite_alignments"`
}

type RiskImpactEvaluationDTO struct {
	RiskMatrixImpactGuidelineID uint                                `json:"risk_matrix_impact_guideline_id"`
	RiskImpactCategoryID        uint                                `json:"risk_impact_category_id"`
	RiskMatrixImpactID          uint                                `json:"risk_matrix_impact_id"`
	RiskMatrixImpactGuideline   RiskMatrixImpactGuidelineDetailsDTO `json:"risk_matrix_impact_guideline"`
	RiskImpactCategory          RiskImpactCategoryDTO               `json:"risk_impact_category"`
	RiskMatrixImpact            RiskMatrixImpactDTO                 `json:"risk_matrix_impact"`
}

type RiskAppetiteAlignmentDTO struct {
	RiskAppetiteStatementID      uint                            `json:"risk_appetite_statement_id"`
	RiskAppetiteToleranceLevelID uint                            `json:"risk_appetite_tolerance_level_id"`
	RiskAppetiteStatement        RiskAppetiteStatementDetailsDTO `json:"risk_appetite_statement"`
	RiskAppetiteToleranceLevel   RiskAppetiteToleranceLevelDTO   `json:"risk_appetite_tolerance_level"`
}

type RiskInherentEvaluationFormDTO struct {
	InherentRiskMatrixLikelihoodID uint                            `json:"inherent_risk_matrix_likelihood_id" binding:"required"`
	InherentRiskMatrixImpactID     uint                            `json:"inherent_risk_matrix_impact_id" binding:"required"`
	RiskImpactEvaluations          []*RiskImpactEvaluationFormDTO  `json:"risk_impact_evaluations"`
	RiskAppetiteAlignments         []*RiskAppetiteAlignmentFormDTO `json:"risk_appetite_alignments"`
}

type RiskImpactEvaluationFormDTO struct {
	RiskMatrixImpactGuidelineID uint `json:"risk_matrix_impact_guideline_id"`
	RiskImpactCategoryID        uint `json:"risk_impact_category_id"`
	RiskMatrixImpactID          uint `json:"risk_matrix_impact_id"`
}

type RiskAppetiteAlignmentFormDTO struct {
	RiskAppetiteStatementID      uint `json:"risk_appetite_statement_id"`
	RiskAppetiteToleranceLevelID uint `json:"risk_appetite_tolerance_level_id"`
}
