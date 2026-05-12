package dtos

type CommitteeTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

type CommitteeTypeListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Type    string
}

type CommitteeTypeFormDTO struct {
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
	Type    string `json:"type" binding:"omitempty,oneof=statutory_reference no_statutory_reference other"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
}
