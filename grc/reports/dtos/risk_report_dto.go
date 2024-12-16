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
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           uint   `json:"count"`
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

type RiskByCategoryDTO struct {
	RiskCategoryID      uint   `json:"risk_category_id"`
	RiskCategoryTitleEn string `json:"risk_category_title_en"`
	RiskCategoryTitleAr string `json:"risk_category_title_ar"`
	Count               uint   `json:"count"`
}

type RiskByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn   string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode string `json:"responsible_unit_color_code"`
	Count                    uint   `json:"count"`
}

type RiskByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           uint   `json:"count"`
}

type RiskByRiskSourceDTO struct {
	RiskSourceID      uint   `json:"risk_source_id"`
	RiskSourceTitleEn string `json:"risk_source_title_en"`
	RiskSourceTitleAr string `json:"risk_source_title_ar"`
	Count             uint   `json:"count"`
}

type RiskByRiskRatingDTO struct {
	RiskRatingID      uint   `json:"risk_rating_id"`
	RiskRatingTitleEn string `json:"risk_rating_title_en"`
	RiskRatingTitleAr string `json:"risk_rating_title_ar"`
	ColorCode         string `json:"color_code"`
	Count             uint   `json:"count"`
}

type RiskByStatusDTO struct {
	RiskStatusID      uint   `json:"risk_status_id"`
	RiskStatusTitleEn string `json:"risk_status_title_en"`
	RiskStatusTitleAr string `json:"risk_status_title_ar"`
	ColorCode         string `json:"color_code"`
	Count             uint   `json:"count"`
}

type RiskByResponseStrategyDTO struct {
	ResponseStrategyID      uint   `json:"response_strategy_id"`
	ResponseStrategyTitleEn string `json:"response_strategy_title_en"`
	ResponseStrategyTitleAr string `json:"response_strategy_title_ar"`
	Count                   uint   `json:"count"`
}

type RiskByRiskAppetiteToleranceLevelDTO struct {
	ToleranceLevelID      uint   `json:"tolerance_level_id"`
	ToleranceLevelTitleEn string `json:"tolerance_level_title_en"`
	ToleranceLevelTitleAr string `json:"tolerance_level_title_ar"`
	ColorCode             string `json:"color_code"`
	Count                 uint   `json:"count"`
}

type RiskByKRIEnabledStatusDTO struct {
	KRIStatus string `json:"kri_status"`
	Count     uint   `json:"count"`
}
type RiskActionPlansByStatusDTO struct {
	ActionPlanStatusID        uint   `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn   string `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   string `json:"action_plan_status_title_ar"`
	ActionPlanStatusType      string `json:"action_plan_status_type"`
	ActionPlanStatusColorCode string `json:"action_plan_status_color_code"`
	Count                     uint   `json:"count"`
}

type RiskActionPlansByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           uint   `json:"count"`
}

type RiskActionPlansBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"severity_level_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelType      string `json:"severity_level_type"`
	Count                  uint   `json:"count"`
}
