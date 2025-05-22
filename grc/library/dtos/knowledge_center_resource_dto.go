package dtos

type KnowledgeCenterResourceListDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	Icon          string  `json:"icon"`
	ModuleTitleEn string  `json:"module_title_en"`
	ModuleTitleAr string  `json:"module_title_ar"`
	StatusTitleEn string  `json:"status_title_en"`
	StatusTitleAr string  `json:"status_title_ar"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}
