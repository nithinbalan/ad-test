package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceIncidentListDTO struct {
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

type ComplianceStrategicObjectiveListDTO struct {
	StrategicObjectiveID      uint   `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr string `json:"strategic_objective_title_ar"`
	StrategicObjectiveStatus  string `json:"strategic_objective_status"`
}

type ComplianceAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}

type ComplianceAssociatedItemControlListResponseDTO struct {
	Items      []*commonDTOs.AssociatedControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ComplianceAssociatedItemIncidentListResponseDTO struct {
	Items      []*ComplianceIncidentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ComplianceAssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*ComplianceStrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ComplianceAssociatedItemsCountResponseDTO struct {
	ComplianceRequirements int `json:"compliance_requirements"`
	Controls               int `json:"controls"`
	Risks                  int `json:"risks"`
	Incidents              int `json:"incidents"`
}

type ComplianceRequirementStatusCountResponseDTO struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type ComplianceRegisterControlMapDTO struct {
	ControlIDs []uint `json:"control_ids" binding:"required"`
}
