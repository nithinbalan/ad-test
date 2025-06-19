package dtos

type AIModuleFeatureListDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Order         int    `json:"order"`

	StatusID      uint   `json:"status_id"`
	StatusTitleEn string `json:"status_title_en"`
	StatusTitleAr string `json:"status_title_ar"`
	Status        string `json:"status"`
	ColorCode     string `json:"color_code"`
}
