package dtos

type ControlEffectivenessRatingDTO struct {
	ID               uint    `json:"id"`
	TitleEn          string  `json:"title_en"`
	TitleAr          string  `json:"title_ar"`
	Type             string  `json:"type"`
	MinScore         float64 `json:"min_score"`
	MaxScore         float64 `json:"max_score"`
	ReductionPercent float64 `json:"reduction_percent"`
	Status           string  `json:"status"`
}
