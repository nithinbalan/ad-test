package dtos

type AssetCategoryDTO struct {
	ID               int    `json:"id"`
	TitleEn          string `json:"title_en"`
	TitleAr          string `json:"title_ar"`
	AssetTypeTitleEn string `json:"asset_type_title_en"`
	AssetTypeTitleAr string `json:"asset_type_title_ar"`
	Status           string `json:"status"`
}

type AssetCategoryFormDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type AssetCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
