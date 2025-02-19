package dtos

type AssetSubCategoryDTO struct {
	ID                   int    `json:"id"`
	AssetCategoryID      *uint  `json:"asset_category_id"`
	TitleEn              string `json:"title_en"`
	TitleAr              string `json:"title_ar"`
	AssetCategoryTitleEn string `json:"asset_category_title_en"`
	AssetCategoryTitleAr string `json:"asset_category_title_ar"`
	Status               string `json:"status"`
}

type AssetSubCategoryFormDTO struct {
	ID              uint   `json:"id"`
	AssetCategoryID uint   `json:"asset_category_id" binding:"required"`
	TitleEn         string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr         string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status          string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type AssetSubCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
