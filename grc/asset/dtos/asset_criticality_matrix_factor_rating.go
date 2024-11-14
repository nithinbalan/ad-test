package dtos

type AssetCriticalityMatrixFactorRatingDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Score     int    `json:"score"`
	Type      string `json:"type,omitempty"`
	ColorCode string `json:"color_code,omitempty"`
}
