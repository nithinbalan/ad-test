package dtos

type CommitteeClassificationListDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	Status        string  `json:"status"`
}

type CommitteeClassificationListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
}

type CommitteeClassificationFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	Status        string  `json:"status" binding:"required,oneof=active inactive"`
}

type CommitteeClassificationDetailsDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	Status        string  `json:"status"`
}
