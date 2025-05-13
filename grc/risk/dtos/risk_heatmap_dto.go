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

type RiskHeatmapResponseByItemDetailsDTO struct {
	RiskMatrixRiskRatingID      uint   `json:"risk_matrix_risk_rating_id"`
	RiskMatrixRiskRatingTitleEN string `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAR string `json:"risk_matrix_risk_rating_title_ar"`
	ColorCode                   string `json:"color_code"`
	RiskRegisterCount           int    `json:"risk_register_count"`
	Value                       uint   `json:"value"`
}

type RiskHeatmapResponseByItemDTO struct {
	ItemID       uint                                   `json:"item_id"`
	ItemTitleEn  string                                 `json:"item_title_en"`
	ItemTitleAr  string                                 `json:"item_title_ar"`
	HeatmapCells []*RiskHeatmapResponseByItemDetailsDTO `json:"heatmap_cells"`
}

type RiskHeatmapQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
