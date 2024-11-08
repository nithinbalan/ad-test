package dtos

type AssetCriticalityMatrixRatingScaleDTO struct {
	ID            uint    `json:"id"`
	Type          string  `json:"type"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en,omitempty"`
	DescriptionAr *string `json:"description_ar,omitempty"`
	ScoreFrom     float64 `json:"score_from"`
	ScoreTo       float64 `json:"score_to"`
	ColorCode     string  `json:"color_code"`
}
