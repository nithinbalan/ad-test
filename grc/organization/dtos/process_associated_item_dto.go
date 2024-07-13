package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ProcessControlListDTO struct {
	ID                                   uint    `json:"id"`
	ProcessID                            uint    `json:"process_id"`
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

type ProcessRiskListDTO struct {
	ID                            uint    `json:"id"`
	RiskRegisterID                uint    `json:"risk_register_id"`
	ProcessID                     uint    `json:"process_id"`
	ReferenceNumber               string  `json:"reference_number"`
	RiskCategoryTitleEn           string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr           string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn             string  `json:"risk_source_title_en"`
	RiskSourceTitleAr             string  `json:"risk_source_title_ar"`
	OwnerUnitTitleEn              string  `json:"owner_unit_title_en"`
	OwnerUnitTitleAr              string  `json:"owner_unit_title_ar"`
	OwnerFirstNameEn              string  `json:"owner_first_name_en"`
	OwnerFirstNameAr              string  `json:"owner_first_name_ar"`
	OwnerLastNameEn               string  `json:"owner_last_name_en"`
	OwnerLastNameAr               string  `json:"owner_last_name_ar"`
	OwnerInitial                  string  `json:"owner_initial"`
	OwnerColor                    string  `json:"owner_color"`
	OwnerEmail                    string  `json:"owner_email"`
	OwnerProfileImagePath         string  `json:"owner_profile_image_path"`
	RiskStatusTitleEn             string  `json:"risk_status_title_en"`
	RiskStatusTitleAr             string  `json:"risk_status_title_ar"`
	RiskStatusColorCode           string  `json:"risk_status_color_code"`
	Title                         string  `json:"title"`
	Description                   *string `json:"description"`
	RiskMatrixLikelihoodTitleEn   string  `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAr   string  `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixLikelihoodColorCode string  `json:"risk_matrix_likelihood_score"`
	RiskMatrixImpactTitleEn       string  `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr       string  `json:"risk_matrix_impact_title_ar"`
	RiskMatrixImpactColorCode     string  `json:"risk_matrix_impact_score"`
	RiskMatrixRiskRatingTitleEn   string  `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr   string  `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingType      string  `json:"risk_matrix_risk_rating_type"`
	RiskMatrixRiskRatingColorCode string  `json:"risk_matrix_risk_rating_color_code"`
}

type ProcessIncidentListDTO struct {
	ID                         uint       `json:"id"`
	IncidentRegisterID         uint       `json:"incident_register_id"`
	ProcessID                  uint       `json:"process_id"`
	ReferenceNumber            string     `json:"reference_number"`
	Title                      string     `json:"title"`
	Description                *string    `json:"description"`
	IncidentAt                 time.Time  `json:"incident_at"`
	Location                   *string    `json:"location"`
	ReportedAt                 *time.Time `json:"reported_at"`
	QuickAction                *string    `json:"quick_action"`
	IncidentCategoryTitleEn    string     `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string     `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn string     `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr string     `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn        string     `json:"incident_type_title_en"`
	IncidentTypeTitleAr        string     `json:"incident_type_title_ar"`
	SeverityLevelTitleEn       string     `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string     `json:"severity_level_title_ar"`
	SeverityLevelColorCode     string     `json:"severity_level_color_code"`
	IncidentStatusTitleEn      string     `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string     `json:"incident_status_title_ar"`
	IncidentStatusColorCode    string     `json:"incident_status_color_code"`
}

type ProcessAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type ProcessAssociatedItemIncidentListResponseDTO struct {
	Items      []*ProcessIncidentListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProcessAssociatedItemRiskListResponseDTO struct {
	Items      []*ProcessRiskListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
