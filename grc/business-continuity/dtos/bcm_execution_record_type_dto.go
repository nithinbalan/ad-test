package dtos

type BCMExecutionRecordTypeDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type BCMExecutionRecordTypeFormDTO struct {
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type BCMExecutionRecordTypeListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
}
