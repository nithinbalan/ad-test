package dtos

type IntegrationCategoryListDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Status        string `json:"status"`
}
