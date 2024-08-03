package dtos

type RiskHeatmapResponseDTO struct {
	RiskMatrixImpactID          uint                        `json:"risk_matrix_impact_id"`
	RiskMatrixImpactTitleEN     string                      `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAR     string                      `json:"risk_matrix_impact_title_ar"`
	RiskMatrixLikelihoodID      uint                        `json:"risk_matrix_likelihood_id"`
	RiskMatrixLikelihoodTitleEN string                      `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAR string                      `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixScore             uint                        `json:"risk_matrix_score"`
	RiskMatrixRiskRating        RiskMatrixRiskRatingDataDTO `json:"risk_matrix_risk_rating"`
	RiskRegisterCount           int                         `json:"risk_register_count"`
}

type RiskHeatmapResponseByItemDTO struct {
	ItemID       uint                      `json:"item_id"`
	ItemTitleEn  string                    `json:"item_title_en"`
	ItemTitleAr  string                    `json:"item_title_ar"`
	HeatmapCells []*RiskHeatmapResponseDTO `json:"heatmap_cells"`
}
