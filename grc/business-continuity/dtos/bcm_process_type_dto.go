package dtos

type BCMProcessTypeDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type BCMProcessTypeFormDTO struct {
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type BCMProcessTypeListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
}
