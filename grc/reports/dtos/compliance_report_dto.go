package dtos

type ComplianceByCategoryDTO struct {
	ComplianceCategoryID      uint   `json:"id"`
	ComplianceCategoryTitleEn string `json:"title_en"`
	ComplianceCategoryTitleAr string `json:"title_ar"`
	Count                     int64  `json:"count"`
}

type ComplianceBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"id"`
	SeverityLevelTitleEn   string `json:"title_en"`
	SeverityLevelTitleAr   string `json:"title_ar"`
	SeverityLevelColorCode string `json:"color_code"`
	SeverityLevelType      string `json:"type"`
	Count                  int64  `json:"count"`
}

type ComplianceByStatusDTO struct {
	ComplianceStatusID        uint   `json:"id"`
	ComplianceStatusTitleEn   string `json:"title_en"`
	ComplianceStatusTitleAr   string `json:"title_ar"`
	ComplianceStatusColorCode string `json:"color_code"`
	ComplianceStatusType      string `json:"type"`
	Count                     int64  `json:"count"`
}

type ComplianceByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    int64  `json:"count"`
}

type ComplianceByResponsibleUserDTO struct {
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
	Count                           int64  `json:"count"`
}

type ComplianceRequirementByComplianceStatusDTO struct {
	ComplianceStatusID      uint   `json:"id"`
	ComplianceStatusTitleEn string `json:"title_en"`
	ComplianceStatusTitleAr string `json:"title_ar"`
	Count                   int64  `json:"count"`
}

type ComplianceRequirementBySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"id"`
	SeverityLevelTitleEn   string `json:"title_en"`
	SeverityLevelTitleAr   string `json:"title_ar"`
	SeverityLevelColorCode string `json:"color_code"`
	SeverityLevelType      string `json:"type"`
	Count                  int64  `json:"count"`
}

type ComplianceRequirementByCategoryDTO struct {
	ComplianceRequirementCategoryID      uint   `json:"id"`
	ComplianceRequirementCategoryTitleEn string `json:"title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"title_ar"`
	Count                                int64  `json:"count"`
}

type ComplianceRequirementByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type ComplianceActionPlanByStatusDTO struct {
	ActionPlanStatusID        uint   `json:"id"`
	ActionPlanStatusTitleEn   string `json:"title_en"`
	ActionPlanStatusTitleAr   string `json:"title_ar"`
	ActionPlanStatusColorCode string `json:"color_code"`
	ActionPlanStatusType      string `json:"type"`
	Count                     int64  `json:"count"`
}

type ComplianceActionPlanByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type SLAByCategoryDTO struct {
	SLACategoryID      uint   `json:"id"`
	SLACategoryTitleEn string `json:"title_en"`
	SLACategoryTitleAr string `json:"title_ar"`
	Count              int64  `json:"count"`
}

type SLABySeverityLevelDTO struct {
	SeverityLevelID        uint   `json:"id"`
	SeverityLevelTitleEn   string `json:"title_en"`
	SeverityLevelTitleAr   string `json:"title_ar"`
	SeverityLevelColorCode string `json:"color_code"`
	SeverityLevelType      string `json:"type"`
	Count                  int64  `json:"count"`
}

type SLAByComplianceStatusDTO struct {
	ComplianceStatusID      uint   `json:"id"`
	ComplianceStatusTitleEn string `json:"title_en"`
	ComplianceStatusTitleAr string `json:"title_ar"`
	Count                   int64  `json:"count"`
}

type ReportListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
