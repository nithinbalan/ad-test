package dtos

type PotentialRiskListDTO struct {
	ID                         uint    `json:"id"`
	RiskCategoryTitleEn        string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string  `json:"risk_category_title_ar"`
	SeverityLevelTitleEn       string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string  `json:"severity_level_title_ar"`
	PotentialRiskStatusTitleEn string  `json:"potential_risk_status_title_en"`
	PotentialRiskStatusTitleAr string  `json:"potential_risk_status_title_ar"`
	Title                      string  `json:"title"`
	Description                *string `json:"description"`
	RiskSource                 string  `json:"risk_source"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  string  `json:"updated_at"`
	CreatedUserFirstNameEn     string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn     *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr     *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn      *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr      *string `json:"updated_user_last_name_ar"`
}
