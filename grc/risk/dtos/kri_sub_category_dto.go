package dtos

type KRISubCategoryDTO struct {
	ID                 int     `json:"id"`
	KRICategoryID      uint    `json:"kri_category_id"`
	TitleEn            string  `json:"title_en"`
	TitleAr            string  `json:"title_ar"`
	DescriptionEn      *string `json:"description_en"`
	DescriptionAr      *string `json:"description_ar"`
	KRICategoryTitleEn string  `json:"kri_category_title_en"`
	KRICategoryTitleAr string  `json:"kri_category_title_ar"`
	Status             string  `json:"status"`
}

type KRISubCategoryFormDTO struct {
	ID            uint    `json:"id"`
	KRICategoryID uint    `json:"kri_category_id" binding:"required"`
	TitleEn       string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type KRISubCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
