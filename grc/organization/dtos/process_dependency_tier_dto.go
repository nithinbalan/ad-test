package dtos

type ProcessDependencyTierDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type ProcessDependencyTierListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
}

type ProcessDependencyTierFormDTO struct {
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
}
