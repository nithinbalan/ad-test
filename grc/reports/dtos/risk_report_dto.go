package dtos

type KRIByCategoryDTO struct {
	KRICategoryID      uint   `json:"id"`
	KRICategoryTitleEn string `json:"title_en"`
	KRICategoryTitleAr string `json:"title_ar"`
	Count              uint   `json:"count"`
}

type KRIByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type KRIByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           uint   `json:"count"`
}

type KRIByRiskLevelDTO struct {
	RiskLevelID      uint   `json:"id"`
	RiskLevelTitleEn string `json:"title_en"`
	RiskLevelTitleAr string `json:"title_ar"`
	ColorCode        string `json:"color_code"`
	Count            uint   `json:"count"`
}

type KRIByReviewFrequencyDTO struct {
	ReviewFrequencyID      uint   `json:"id"`
	ReviewFrequencyTitleEn string `json:"title_en"`
	ReviewFrequencyTitleAr string `json:"title_ar"`
	Count                  uint   `json:"count"`
}

type RiskByCategoryDTO struct {
	RiskCategoryID      uint   `json:"id"`
	RiskCategoryTitleEn string `json:"title_en"`
	RiskCategoryTitleAr string `json:"title_ar"`
	Count               uint   `json:"count"`
}

type RiskByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type RiskByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           uint   `json:"count"`
}

type RiskByRiskSourceDTO struct {
	RiskSourceID      uint   `json:"id"`
	RiskSourceTitleEn string `json:"title_en"`
	RiskSourceTitleAr string `json:"title_ar"`
	Count             uint   `json:"count"`
}

type RiskByRiskRatingDTO struct {
	RiskRatingID      uint   `json:"id"`
	RiskRatingTitleEn string `json:"title_en"`
	RiskRatingTitleAr string `json:"title_ar"`
	ColorCode         string `json:"color_code"`
	Count             uint   `json:"count"`
}

type RiskByStatusDTO struct {
	RiskStatusID      uint   `json:"id"`
	RiskStatusTitleEn string `json:"title_en"`
	RiskStatusTitleAr string `json:"title_ar"`
	ColorCode         string `json:"color_code"`
	Count             uint   `json:"count"`
}

type RiskByResponseStrategyDTO struct {
	ResponseStrategyID      uint   `json:"id"`
	ResponseStrategyTitleEn string `json:"title_en"`
	ResponseStrategyTitleAr string `json:"title_ar"`
	Count                   uint   `json:"count"`
}

type RiskByRiskAppetiteToleranceLevelDTO struct {
	ToleranceLevelID      uint   `json:"id"`
	ToleranceLevelTitleEn string `json:"title_en"`
	ToleranceLevelTitleAr string `json:"title_ar"`
	ColorCode             string `json:"color_code"`
	Count                 uint   `json:"count"`
}

type RiskByKRIEnabledStatusDTO struct {
	KRIStatus string `json:"kri_status"`
	Count     uint   `json:"count"`
}

type RiskActionPlansByStatusDTO struct {
	ActionPlanStatusID        uint   `json:"id"`
	ActionPlanStatusTitleEn   string `json:"title_en"`
	ActionPlanStatusTitleAr   string `json:"title_ar"`
	ActionPlanStatusType      string `json:"type"`
	ActionPlanStatusColorCode string `json:"color_code"`
	Count                     uint   `json:"count"`
}

type RiskActionPlansByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           uint   `json:"count"`
}

type RiskActionPlansBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"id"`
	SeverityLevelTitleEn   string `json:"title_en"`
	SeverityLevelTitleAr   string `json:"title_ar"`
	SeverityLevelColorCode string `json:"color_code"`
	SeverityLevelType      string `json:"type"`
	Count                  uint   `json:"count"`
}
