package dtos

type IncidentStatusDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	ColorCode     string `json:"color_code"`
	Type          string `json:"type"`
	Status        string `json:"status"`
}
