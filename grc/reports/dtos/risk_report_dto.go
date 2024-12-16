package dtos

type KRIByCategoryDTO struct {
	KRICategoryID      uint   `json:"kri_category_id"`
	KRICategoryTitleEn string `json:"kri_category_title_en"`
	KRICategoryTitleAr string `json:"kri_category_title_ar"`
	Count              uint   `json:"count"`
}

type KRIByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn   string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode string `json:"responsible_unit_color_code"`
	Count                    uint   `json:"count"`
}

type KRIByResponsibleUserDTO struct {
	ResponsibleUserID          uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn  string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr  string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial     string `json:"responsible_user_initial"`
	ResponsibleUserColor       string `json:"responsible_user_color"`
	ResponsibleUserEmail       string `json:"responsible_user_email"`
	Count                      uint   `json:"count"`
}

type KRIByRiskLevelDTO struct {
	RiskLevelID      uint   `json:"risk_level_id"`
	RiskLevelTitleEn string `json:"risk_level_title_en"`
	RiskLevelTitleAr string `json:"risk_level_title_ar"`
	ColorCode        string `json:"color_code"`
	Count            uint   `json:"count"`
}

type KRIByReviewFrequencyDTO struct {
	ReviewFrequencyID      uint   `json:"review_frequency_id"`
	ReviewFrequencyTitleEn string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr string `json:"review_frequency_title_ar"`
	Count                  uint   `json:"count"`
}
