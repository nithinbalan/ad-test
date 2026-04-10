package dtos

type RiskSubCategoryDTO struct {
	ID                  int     `json:"id"`
	RiskCategoryID      uint    `json:"risk_category_id"`
	TitleEn             string  `json:"title_en"`
	TitleAr             string  `json:"title_ar"`
	DescriptionEn       *string `json:"description_en"`
	DescriptionAr       *string `json:"description_ar"`
	RiskCategoryTitleEn string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr string  `json:"risk_category_title_ar"`
	Status              string  `json:"status"`
}

type RiskSubCategoryFormDTO struct {
	ID             uint    `json:"id"`
	RiskCategoryID uint    `json:"risk_category_id" binding:"required"`
	TitleEn        string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr        string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn  *string `json:"description_en"`
	DescriptionAr  *string `json:"description_ar"`
	Status         string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type RiskSubCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}