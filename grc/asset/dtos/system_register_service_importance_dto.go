package dtos

type SystemRegisterServiceImportanceDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type SystemRegisterServiceImportanceFormDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type SystemRegisterServiceImportanceListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
