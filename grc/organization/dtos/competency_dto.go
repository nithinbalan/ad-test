package dtos

type CompetencyWithGroupDTO struct {
	ID                     uint   `json:"id"`
	TitleEn                string `json:"title_en" `
	TitleAr                string `json:"title_ar" `
	Status                 string `json:"status" `
	CompetencyGroupID      uint   `json:"competency_group_id"`
	CompetencyGroupTitleEn string `json:"competency_group_title_en"`
	CompetencyGroupTitleAr string `json:"competency_group_title_ar"`
}

type CompetencyFormDTO struct {
	CompetencyGroupID uint   `json:"competency_group_id" binding:"required"`
	TitleEn           string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr           string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	Status            string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type CompetencyListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
