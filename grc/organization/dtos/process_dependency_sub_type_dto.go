package dtos

type ProcessDependencySubTypeDTO struct {
	ID                           int     `json:"id"`
	ProcessDependencyTypeID      *uint   `json:"process_dependency_type_id"`
	TitleEn                      string  `json:"title_en"`
	TitleAr                      string  `json:"title_ar"`
	DescriptionEn                *string `json:"description_en"`
	DescriptionAr                *string `json:"description_ar"`
	ProcessDependencyTypeTitleEn string  `json:"process_dependency_type_title_en"`
	ProcessDependencyTypeTitleAr string  `json:"process_dependency_type_title_ar"`
	Status                       string  `json:"status"`
}

type ProcessDependencySubTypeFormDTO struct {
	ID                      uint    `json:"id"`
	ProcessDependencyTypeID uint    `json:"process_dependency_type_id" binding:"required"`
	TitleEn                 string  `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr                 string  `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	DescriptionEn           *string `json:"description_en" gorm:"type:nvarchar(MAX)"`
	DescriptionAr           *string `json:"description_ar" gorm:"type:nvarchar(MAX)"`
	Status                  string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type ProcessDependencySubTypeListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
