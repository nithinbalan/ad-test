package dtos

type IncidentSubCategoryDTO struct {
	ID                      int    `json:"id"`
	IncidentCategoryID      uint   `json:"incident_category_id"`
	TitleEn                 string `json:"title_en"`
	TitleAr                 string `json:"title_ar"`
	IncidentCategoryTitleEn string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr string `json:"incident_category_title_ar"`
	Status                  string `json:"status"`
}
type IncidentSubCategoryFormDTO struct {
	ID                 uint   `json:"id"`
	IncidentCategoryID uint   `json:"incident_category_id" binding:"required"`
	TitleEn            string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr            string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status             string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type IncidentSubCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
