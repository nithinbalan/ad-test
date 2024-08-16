package dtos

import organizationDTOs "jethur/grc/organization/dtos"

type RiskControlAnalysisListDTO struct {
	ID                                   uint    `json:"id"`
	ControlID                            uint    `json:"control_id"`
	ControlTestID                        *uint   `json:"control_test_id"`
	ControlReferenceNumber               string  `json:"control_reference_number"`
	ControlCode                          *string `json:"control_code"`
	ControlTitleEn                       string  `json:"control_title_en"`
	ControlTitleAr                       string  `json:"control_title_ar"`
	Score                                string  `json:"score"`
	ControlDescription                   *string `json:"control_description"`
	ControlTypeTitleEn                   string  `json:"control_type_title_en"`
	ControlTypeTitleAr                   string  `json:"control_type_title_ar"`
	ControlCategoryTitleEn               string  `json:"control_category_title_en"`
	ControlCategoryTitleAr               string  `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string `json:"owner_unit_title_ar"`
	OwnerFirstNameEn                     *string `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string `json:"owner_last_name_ar"`
	Type                                 string  `json:"type" binding:"required"`
	ControlImplementationStatusTitleEn   *string `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   *string `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode *string `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string `json:"control_efficiency_status_color_code"`
	CreatedAt                            string  `json:"created_at"`
	CreatedUserID                        uint    `json:"created_user_id"`
	CreatedUserFirstNameEn               string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string  `json:"created_user_last_name_ar"`
}

type RiskControlAnalysisListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type RiskControlAnalysisFormDTO struct {
	ControlEfficiencyStatusID     uint `json:"control_efficiency_status_id" binding:"required"`
	ControlImplementationStatusID uint `json:"control_implementation_status_id" binding:"required"`
}

type RiskControlAnalysisDetailDTO struct {
	ID                          uint                                            `json:"id"`
	ControlTestID               *uint                                           `json:"control_test_id"`
	ControlEfficiencyStatus     organizationDTOs.ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	ControlImplementationStatus organizationDTOs.ControlImplementationStatusDTO `json:"control_implementation_status" `
	Score                       uint                                            `json:"score"`
}

type RiskResidualControlAnalysisDetailDTO struct {
	ID                          uint                                            `json:"id"`
	ControlTestID               *uint                                           `json:"control_test_id"`
	ControlEfficiencyStatus     organizationDTOs.ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	ControlImplementationStatus organizationDTOs.ControlImplementationStatusDTO `json:"control_implementation_status" `
	Score                       uint                                            `json:"score"`
}
