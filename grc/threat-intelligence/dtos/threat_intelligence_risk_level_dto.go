package dtos

type ThreatIntelligenceRiskLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type ThreatIntelligenceRiskLevelFormDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Type      string `json:"type" binding:"required,max=10,oneof=very_low low medium high very_high" gorm:"type:varchar(10);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type ThreatIntelligenceRiskLevelListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
