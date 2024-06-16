package dtos

import commonDTOs "jethur/grc/common/dtos"

type IncidentControlListDTO struct {
	ControlID                            uint    `json:"control_id"`
	ControlReferenceNumber               string  `json:"control_reference_number"`
	ControlCode                          *string `json:"control_code"`
	ControlLastTestDate                  *string `json:"control_last_test_date"`
	ControlCategoryTitleEn               string  `json:"control_category_title_en"`
	ControlCategoryTitleAr               string  `json:"control_category_title_ar"`
	ControlTypeTitleEn                   string  `json:"control_type_title_en"`
	ControlTypeTitleAr                   string  `json:"control_type_title_ar"`
	ControlTitleEn                       string  `json:"control_title_en"`
	ControlTitleAr                       string  `json:"control_title_ar"`
	ControlDescription                   *string `json:"control_description"`
	ControlEfficiencyStatusEn            string  `json:"control_efficiency_status_en"`
	ControlEfficiencyStatusAr            string  `json:"control_efficiency_status_ar"`
	ControlEfficiencyStatusType          string  `json:"control_efficiency_status_type"`
	ControlEfficiencyStatusColorCode     string  `json:"control_efficiency_status_color_code"`
	ControlImplementationStatusEn        string  `json:"control_implementation_status_en"`
	ControlImplementationStatusAr        string  `json:"control_implementation_status_ar"`
	ControlImplementationStatusType      string  `json:"control_implementation_status_type"`
	ControlImplementationStatusColorCode string  `json:"control_implementation_status_color_code"`
	ControlOwnerUnitTitleEn              *string `json:"control_owner_unit_title_en"`
	ControlOwnerUnitTitleAr              *string `json:"control_owner_unit_title_ar"`
	ControlOwnerFirstNameEn              *string `json:"control_owner_first_name_en"`
	ControlOwnerFirstNameAr              *string `json:"control_owner_first_name_ar"`
	ControlOwnerLastNameEn               *string `json:"control_owner_last_name_en"`
	ControlOwnerLastNameAr               *string `json:"control_owner_last_name_ar"`
	ControlOwnerInitial                  string  `json:"control_owner_initial"`
	ControlOwnerColor                    string  `json:"control_owner_color"`
	ControlOwnerEmail                    string  `json:"control_owner_email"`
	ControlOwnerProfileImagePath         string  `json:"control_owner_profile_image_path"`
}

type IncidentComplianceRequirementListDTO struct {
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitle           string `json:"compliance_requirement_title"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type IncidentRiskListDTO struct {
	ReferenceNumber       string  `json:"reference_number"`
	RiskCategoryTitleEn   string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr   string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn     string  `json:"risk_source_title_en"`
	RiskSourceTitleAr     string  `json:"risk_source_title_ar"`
	RiskStatusTitleEn     string  `json:"risk_status_title_en"`
	RiskStatusTitleAr     string  `json:"risk_status_title_ar"`
	RiskStatusColorCode   string  `json:"risk_status_color_code"`
	Title                 string  `json:"title"`
	Description           *string `json:"description"`
	OwnerUnitTitleEn      string  `json:"owner_unit_title_en"`
	OwnerUnitTitleAr      string  `json:"owner_unit_title_ar"`
	OwnerFirstNameEn      string  `json:"owner_first_name_en"`
	OwnerFirstNameAr      string  `json:"owner_first_name_ar"`
	OwnerLastNameEn       string  `json:"owner_last_name_en"`
	OwnerLastNameAr       string  `json:"owner_last_name_ar"`
	OwnerInitial          string  `json:"owner_initial"`
	OwnerColor            string  `json:"owner_color"`
	OwnerEmail            string  `json:"owner_email"`
	OwnerProfileImagePath string  `json:"owner_profile_image_path"`
}

type IncidentProcessListDTO struct {
	ProcessGroupTitleEn      string  `json:"process_group_title_en"`
	ProcessGroupTitleAr      string  `json:"process_group_title_ar"`
	OwnerUnitTitleEn         *string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr         *string `json:"owner_unit_title_ar"`
	OwnerFirstNameEn         *string `json:"owner_first_name_en"`
	OwnerFirstNameAr         *string `json:"owner_first_name_ar"`
	OwnerLastNameEn          *string `json:"owner_last_name_en"`
	OwnerLastNameAr          *string `json:"owner_last_name_ar"`
	OwnerInitial             string  `json:"owner_initial"`
	OwnerColor               string  `json:"owner_color"`
	OwnerEmail               string  `json:"owner_email"`
	OwnerProfileImagePath    string  `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn   string  `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr   string  `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode string  `json:"lifecycle_status_color_code"`
	ReferenceNumber          string  `json:"reference_number"`
	TitleEn                  string  `json:"title_en"`
	TitleAr                  string  `json:"title_ar"`
	Description              *string `json:"description"`
}

type IncidentAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type IncidentAssociatedItemControlListResponseDTO struct {
	Items      []*IncidentControlListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*IncidentComplianceRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}

type IncidentAssociatedItemRiskListResponseDTO struct {
	Items      []*IncidentRiskListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentAssociatedItemProcessListResponseDTO struct {
	Items      []*IncidentProcessListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
