package dtos

type ProcessEquipmentToolDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type ProcessEquipmentToolListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
}

type ProcessEquipmentToolFormDTO struct {
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
}
