package dtos

type IncidentImpactDTO struct {
	ID          uint    `json:"id"`
	TitleEn     string  `json:"title_en"`
	TitleAr     string  `json:"title_ar"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}
