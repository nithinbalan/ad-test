package dtos

type ProjectImpactAssessmentTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

