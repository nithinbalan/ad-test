package dtos

type ProjectMilestoneStatusDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Type          string `json:"type"`
	ColorCode     string `json:"color_code"`
	Status        string `json:"status"`
}
