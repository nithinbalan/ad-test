package dtos

type AuditFindingCategoryDTO struct {
	ID            int    `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	ColorCode     string `json:"color_code"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
}
