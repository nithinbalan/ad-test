package dtos

type ProjectClassificationDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type ProjectClassificationFormDTO struct {
	TitleEn string `json:"title_en" binding:"required,max=100"`
	TitleAr string `json:"title_ar" binding:"required,max=100"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive"`
}
