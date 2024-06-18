package dtos

type RiskInherentEvaluationDetailsDTO struct {
	ID                        uint                         `json:"id"`
	Type                      string                       `json:"type"`
	RiskMatrixLikelihood      RiskMatrixLikelihoodDTO      `json:"risk_matrix_likelihood"`
	RiskMatrixImpact          RiskMatrixImpactDTO          `json:"risk_matrix_impact"`
	RiskMatrixRiskRating      RiskMatrixRiskRatingDTO      `json:"risk_matrix_risk_rating"`
	RiskMatrixRiskRatingValue RiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating_value"`
	RiskImpactAnalyses        []RiskImpactAnalysisDTO      `json:"risk_impact_analyses"`
	RiskAppetiteAlignments    []RiskAppetiteAlignmentDTO   `json:"risk_appetite_alignments"`
}

type RiskImpactAnalysisDTO struct {
	ID                           uint                             `json:"id"`
	RiskEvaluationID             uint                             `json:"risk_evaluation_id"`
	RiskImpactCategory           RiskImpactCategoryDTO            `json:"risk_impact_category"`
	RiskMatrixImpact             RiskMatrixImpactDTO              `json:"risk_matrix_impact"`
	RiskImpactAnalysisGuidelines []RiskImpactAnalysisGuidelineDTO `json:"risk_impact_analysis_guidelines"`
}

type RiskImpactAnalysisGuidelineDTO struct {
	ID                          uint                                `json:"id"`
	RiskImpactAnalysisID        uint                                `json:"risk_impact_analysis_id"`
	RiskMatrixImpactGuidelineID uint                                `json:"risk_matrix_impact_guideline_id"`
	RiskMatrixImpactGuideline   RiskMatrixImpactGuidelineDetailsDTO `json:"risk_matrix_impact_guideline"`
}

type RiskAppetiteAlignmentDTO struct {
	RiskAppetiteStatementID      uint                            `json:"risk_appetite_statement_id"`
	RiskAppetiteToleranceLevelID uint                            `json:"risk_appetite_tolerance_level_id"`
	RiskAppetiteStatement        RiskAppetiteStatementDetailsDTO `json:"risk_appetite_statement"`
	RiskAppetiteToleranceLevel   RiskAppetiteToleranceLevelDTO   `json:"risk_appetite_tolerance_level"`
}

type RiskInherentEvaluationFormDTO struct {
	RiskMatrixLikelihoodID uint                            `json:"risk_matrix_likelihood_id" binding:"required"`
	RiskMatrixImpactID     uint                            `json:"risk_matrix_impact_id" binding:"required"`
	RiskImpactAnalyses     []*RiskImpactAnalysisFormDTO    `json:"risk_impact_analyses"`
	RiskAppetiteAlignments []*RiskAppetiteAlignmentFormDTO `json:"risk_appetite_alignments"`
}

type RiskImpactAnalysisFormDTO struct {
	RiskMatrixImpactGuidelineIDs []uint `json:"risk_matrix_impact_guideline_ids"`
	RiskImpactCategoryID         uint   `json:"risk_impact_category_id"`
	RiskMatrixImpactID           uint   `json:"risk_matrix_impact_id"`
}

type RiskAppetiteAlignmentFormDTO struct {
	RiskAppetiteStatementID      uint `json:"risk_appetite_statement_id"`
	RiskAppetiteToleranceLevelID uint `json:"risk_appetite_tolerance_level_id"`
}
