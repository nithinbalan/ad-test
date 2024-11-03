package dtos

type ComplianceByCategoryDTO struct {
	ComplianceCategoryID      uint   `json:"compliance_category_id"`
	ComplianceCategoryTitleEn string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr string `json:"compliance_category_title_ar"`
	Count                     int64  `json:"count"`
}

type ComplianceBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"severity_level_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelType      string `json:"severity_level_type"`
	Count                  int64  `json:"count"`
}

type ComplianceByStatusDTO struct {
	ComplianceStatusID        uint   `json:"compliance_status_id"`
	ComplianceStatusTitleEn   string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr   string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode string `json:"compliance_status_color_code"`
	ComplianceStatusType      string `json:"compliance_status_type"`
	Count                     int64  `json:"count"`
}

type ComplianceByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn   string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string `json:"responsible_unit_title_ar"`
	ResponsibleUnitColorCode string `json:"responsible_unit_color_code"`
	Count                    int64  `json:"count"`
}

type ComplianceByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           int64  `json:"count"`
}

type ComplianceRequirementByComplianceStatusDTO struct {
	ComplianceStatusID      uint   `json:"compliance_status_id"`
	ComplianceStatusTitleEn string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr string `json:"compliance_status_title_ar"`
	Count                   int64  `json:"count"`
}

type ComplianceRequirementBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"severity_level_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelType      string `json:"severity_level_type"`
	Count                  int64  `json:"count"`
}

type ComplianceRequirementByCategoryDTO struct {
	ComplianceRequirementCategoryID      uint   `json:"compliance_requirement_category_id"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	Count                                int64  `json:"count"`
}

type ComplianceRequirementByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           int64  `json:"count"`
}

type ComplianceActionPlanByStatusDTO struct {
	ActionPlanStatusID        uint   `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn   string `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   string `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode string `json:"action_plan_status_color_code"`
	ActionPlanStatusType      string `json:"action_plan_status_type"`
	Count                     int64  `json:"count"`
}

type ComplianceActionPlanByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           int64  `json:"count"`
}
type SLAByCategoryDTO struct {
	SLACategoryID      uint   `json:"sla_category_id"`
	SLACategoryTitleEn string `json:"sla_category_title_en"`
	SLACategoryTitleAr string `json:"sla_category_title_ar"`
	Count              int64  `json:"count"`
}
type SLABySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"severity_level_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelType      string `json:"severity_level_type"`
	Count                  int64  `json:"count"`
}
type SLAByComplianceStatusDTO struct {
	ComplianceStatusID      uint   `json:"compliance_status_id"`
	ComplianceStatusTitleEn string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr string `json:"compliance_status_title_ar"`
	Count                   int64  `json:"count"`
}

type ReportListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
