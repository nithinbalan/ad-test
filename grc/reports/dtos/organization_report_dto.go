package dtos

type UsersByStatusDTO struct {
	UserStatusID        uint   `json:"id"`
	UserStatusTitleEn   string `json:"title_en"`
	UserStatusTitleAr   string `json:"title_ar"`
	UserStatusColorCode string `json:"color_code"`
	UserStatusType      string `json:"type"`
	Count               int    `json:"count"`
}

type UsersByRoleDTO struct {
	RoleID      uint   `json:"id"`
	RoleTitleEn string `json:"title_en"`
	RoleTitleAr string `json:"title_ar"`
	Count       int    `json:"count"`
}

type UsersByDesignationDTO struct {
	DesignationID      uint   `json:"id"`
	DesignationTitleEn string `json:"title_en"`
	DesignationTitleAr string `json:"title_ar"`
	Count              int    `json:"count"`
}

type ProcessByUnitDTO struct {
	OwnerUnitID        uint   `json:"id"`
	OwnerUnitTitleEn   string `json:"title_en"`
	OwnerUnitTitleAr   string `json:"title_ar"`
	OwnerUnitColorCode string `json:"color_code"`
	OwnerUnitStatus    string `json:"status"`
	Count              int    `json:"count"`
}

type ProcessByRiskRatingDTO struct {
	RiskRatingID        uint   `json:"id"`
	RiskRatingTitleEn   string `json:"title_en"`
	RiskRatingTitleAr   string `json:"title_ar"`
	RiskRatingColorCode string `json:"color_code"`
	RiskRatingType      string `json:"type"`
	Count               int    `json:"count"`
}

type ProcessByResponsibleUserDTO struct {
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
	Count                           int    `json:"count"`
}

type ProcessByProcessGroupDTO struct {
	ProcessGroupID      uint   `json:"id"`
	ProcessGroupTitleEn string `json:"title_en"`
	ProcessGroupTitleAr string `json:"title_ar"`
	Count               int    `json:"count"`
}

type ControlsByControlTypeDTO struct {
	ControlTypeID      uint   `json:"id"`
	ControlTypeTitleEn string `json:"title_en"`
	ControlTypeTitleAr string `json:"title_ar"`
	Count              int    `json:"count"`
}

type ControlsByControlCategoryDTO struct {
	ControlCategoryID      uint   `json:"id"`
	ControlCategoryTitleEn string `json:"title_en"`
	ControlCategoryTitleAr string `json:"title_ar"`
	Count                  int    `json:"count"`
}

type ControlsByComplianceRegisterDTO struct {
	ComplianceRegisterID      uint   `json:"id"`
	ComplianceRegisterTitleEn string `json:"title_en"`
	ComplianceRegisterTitleAr string `json:"title_ar"`
	Count                     int    `json:"count"`
}

type ControlsByRiskRatingDTO struct {
	RiskRatingID        uint   `json:"id"`
	RiskRatingTitleEn   string `json:"title_en"`
	RiskRatingTitleAr   string `json:"title_ar"`
	RiskRatingColorCode string `json:"color_code"`
	RiskRatingType      string `json:"type"`
	Count               int    `json:"count"`
}

type ControlsByImplementationStatusDTO struct {
	ControlImplementationStatusID       uint   `json:"id"`
	ControlImplementationStatusTitleEn  string `json:"title_en"`
	ControlImplementationStatusTitleAr  string `json:"title_ar"`
	ControlImplementationStatusolorCode string `json:"color_code"`
	ControlImplementationStatusType     string `json:"type"`
	Count                               int    `json:"count"`
}

type ControlsByEfficiencyDTO struct {
	ControlEfficiencyStatusID        uint   `json:"id"`
	ControlEfficiencyStatusTitleEn   string `json:"title_en"`
	ControlEfficiencyStatusTitleAr   string `json:"title_ar"`
	ControlEfficiencyStatusColorCode string `json:"color_code"`
	ControlEfficiencyStatusType      string `json:"type"`
	Count                            int    `json:"count"`
}

type ControlsByApplicabilityDTO struct {
	IsApplicable bool `json:"is_applicable"`
	Count        int  `json:"count"`
}
