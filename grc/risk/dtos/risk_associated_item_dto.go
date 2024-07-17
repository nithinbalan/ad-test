package dtos

import commonDTOs "jethur/grc/common/dtos"

type RiskControlListDTO struct {
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
	ControlEfficiencyStatusTitleEn       string  `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string  `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusType          string  `json:"control_efficiency_status_type"`
	ControlEfficiencyStatusColorCode     string  `json:"control_efficiency_status_color_code"`
	ControlImplementationStatusTitleEn   string  `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string  `json:"control_implementation_status_title_ar"`
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

type RiskComplianceRequirementListDTO struct {
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitle           string `json:"compliance_requirement_title"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type RiskProcessListDTO struct {
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

type RiskStrategicObjectiveListDTO struct {
	StrategicObjectiveID      uint    `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr *string `json:"strategic_objective_title_ar"`
	StrategicObjectiveStatus  *string `json:"strategic_objective_status"`
}
type RiskSLAListDTO struct {
	ReferenceNumber                 string `json:"reference_number"`
	LifecycleStatusID               uint   `json:"lifecycle_status_id"`
	LifecycleStatusTitleEn          string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr          string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode        string `json:"lifecycle_status_color_code"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	Description                     string `json:"description"`
	Client                          string `json:"client"`
	SLACategoryTitleEn              string `json:"sla_category_title_en"`
	SLACategoryTitleAr              string `json:"sla_category_title_ar"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	ComplianceStatusTitleEn         string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string `json:"compliance_status_color_code"`
}

type RiskIncidentListDTO struct {
	IncidentRegisterID              uint   `json:"incident_register_id"`
	IncidentReferenceNumber         string `json:"incident_reference_number"`
	IncidentRegisterTitle           string `json:"incident_register_title"`
	IncidentRegisterDescription     string `json:"incident_register_description"`
	IncidentCategoryTitleEn         string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr         string `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn      string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr      string `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn             string `json:"incident_type_title_en"`
	IncidentTypeTitleAr             string `json:"incident_type_title_ar"`
	IncidentSeverityLevelTitleEn    string `json:"incident_severity_level_title_en"`
	IncidentSeverityLevelTitleAr    string `json:"incident_severity_level_title_ar"`
	IncidentSeverityLevelColorCode  string `json:"incident_severity_level_color_code"`
	IncidentStatusTitleEn           string `json:"incident_status_title_en"`
	IncidentStatusTitleAr           string `json:"incident_status_title_ar"`
	IncidentStatusColorCode         string `json:"incident_status_color_code"`
	IncidentReportedUserFirstNameEn string `json:"incident_reported_user_first_name_en"`
	IncidentReportedUserFirstNameAr string `json:"incident_reported_user_first_name_ar"`
	IncidentReportedUserLastNameEn  string `json:"incident_reported_user_last_name_en"`
	IncidentReportedUserLastNameAr  string `json:"incident_reported_user_last_name_ar"`
}

type RiskAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type RiskAssociatedItemControlListResponseDTO struct {
	Items      []*RiskControlListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*RiskComplianceRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

type RiskAssociatedItemProcessListResponseDTO struct {
	Items      []*RiskProcessListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*RiskStrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type RiskAssociatedItemSLAListResponseDTO struct {
	Items      []*RiskSLAListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemIncidentListResponseDTO struct {
	Items      []*RiskIncidentListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
