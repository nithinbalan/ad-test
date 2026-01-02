package dtos

type ThirdPartyByTypeDTO struct {
	TypeID  uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type ThirdPartyByCountryDTO struct {
	CountryID uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Count     uint   `json:"count"`
}

type ThirdPartyByDepartmentDTO struct {
	DepartmentID uint   `json:"id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Count        uint   `json:"count"`
}

type ThirdPartyByCriticalityDTO struct {
	SeverityID uint   `json:"id"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	ColorCode  string `json:"color_code"`
	Count      uint   `json:"count"`
}

type ThirdPartyByRiskRatingDTO struct {
	RiskLevelID uint   `json:"id"`
	TitleEn     string `json:"title_en"`
	TitleAr     string `json:"title_ar"`
	ColorCode   string `json:"color_code"`
	Count       int64  `json:"count"`
}

type ThirdPartyServiceExpiryDTO struct {
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Category string `json:"category"`
	Count    uint   `json:"count"`
}

type TPAssessmentsByStatusDTO struct {
	StatusID  uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int64  `json:"count"`
}

type TPAssessmentsByRiskRatingDTO struct {
	RiskLevelID uint   `json:"id"`
	TitleEn     string `json:"title_en"`
	TitleAr     string `json:"title_ar"`
	ColorCode   string `json:"color_code"`
	Count       int64  `json:"count"`
}

type TPAssessmentsByCategoryDTO struct {
	CategoryID uint   `json:"id"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	Count      int64  `json:"count"`
}

type TPAssessmentsByTemplateDTO struct {
	TemplateID uint   `json:"id"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	Count      int64  `json:"count"`
}

type TPOverdueAssessmentsDTO struct {
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Category string `json:"category"`
	Count    int64  `json:"count"`
}
