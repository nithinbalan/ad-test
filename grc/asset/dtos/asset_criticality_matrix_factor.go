package dtos

type AssetCriticalityMatrixFactorDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en,omitempty"`
	DescriptionAr *string `json:"description_ar,omitempty"`
	Weightage     float64 `json:"weightage"`
}
