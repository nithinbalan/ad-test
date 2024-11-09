package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AssetControlListDTO struct {
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
	ControlOwnerID                       *uint   `json:"control_owner_id"`
	ControlOwnerFirstNameEn              *string `json:"control_owner_first_name_en"`
	ControlOwnerFirstNameAr              *string `json:"control_owner_first_name_ar"`
	ControlOwnerLastNameEn               *string `json:"control_owner_last_name_en"`
	ControlOwnerLastNameAr               *string `json:"control_owner_last_name_ar"`
	ControlOwnerInitial                  string  `json:"control_owner_initial"`
	ControlOwnerColor                    string  `json:"control_owner_color"`
	ControlOwnerEmail                    string  `json:"control_owner_email"`
	ControlOwnerProfileImagePath         string  `json:"control_owner_profile_image_path"`
}
type AssetIncidentListDTO struct {
	IncidentRegisterID                   uint   `json:"incident_register_id"`
	IncidentReferenceNumber              string `json:"incident_reference_number"`
	IncidentRegisterTitle                string `json:"incident_register_title"`
	IncidentRegisterDescription          string `json:"incident_register_description"`
	IncidentCategoryTitleEn              string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr              string `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn           string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr           string `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn                  string `json:"incident_type_title_en"`
	IncidentTypeTitleAr                  string `json:"incident_type_title_ar"`
	IncidentSeverityLevelTitleEn         string `json:"incident_severity_level_title_en"`
	IncidentSeverityLevelTitleAr         string `json:"incident_severity_level_title_ar"`
	IncidentSeverityLevelColorCode       string `json:"incident_severity_level_color_code"`
	IncidentStatusTitleEn                string `json:"incident_status_title_en"`
	IncidentStatusTitleAr                string `json:"incident_status_title_ar"`
	IncidentStatusColorCode              string `json:"incident_status_color_code"`
	IncidentReportedUserID               uint   `json:"incident_reported_user_id"`
	IncidentReportedUserFirstNameEn      string `json:"incident_reported_user_first_name_en"`
	IncidentReportedUserFirstNameAr      string `json:"incident_reported_user_first_name_ar"`
	IncidentReportedUserLastNameEn       string `json:"incident_reported_user_last_name_en"`
	IncidentReportedUserLastNameAr       string `json:"incident_reported_user_last_name_ar"`
	IncidentReportedUserInitial          string `json:"incident_reported_user_initial"`
	IncidentReportedUserColor            string `json:"incident_reported_user_color"`
	IncidentReportedUserEmail            string `json:"incident_reported_user_email"`
	IncidentReportedUserProfileImagePath string `json:"incident_reported_user_profile_image_path"`
}

type AssetAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type AssetAssociatedItemControlListResponseDTO struct {
	Items      []*AssetControlListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetAssociatedItemIncidentListResponseDTO struct {
	Items      []*AssetIncidentListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
