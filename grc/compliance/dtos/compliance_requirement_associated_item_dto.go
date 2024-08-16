package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ComplianceRequirementControlListDTO struct {
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

type ComplianceRequirementStrategicObjectiveListDTO struct {
	StrategicObjectiveID      uint    `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr *string `json:"strategic_objective_title_ar"`
	StrategicObjectiveStatus  *string `json:"strategic_objective_status"`
}

type ComplianceRequirementRiskListDTO struct {
	RiskRegisterID                uint       `json:"risk_register_id"`
	RiskReferenceNumber           string     `json:"risk_reference_number"`
	RiskRegisterTitle             string     `json:"risk_register_title"`
	RiskRegisterDescription       string     `json:"risk_register_description"`
	RiskCategoryTitleEn           string     `json:"risk_category_title_en"`
	RiskCategoryTitleAr           string     `json:"risk_category_title_ar"`
	RiskSourceTitleEn             string     `json:"risk_source_title_en"`
	RiskSourceTitleAr             string     `json:"risk_source_title_ar"`
	RiskOwnerUnitTitleEn          string     `json:"risk_owner_unit_title_en"`
	RiskOwnerUnitTitleAr          string     `json:"risk_owner_unit_title_ar"`
	RiskOwnerID                   uint       `json:"risk_owner_id"`
	RiskOwnerFirstNameEn          string     `json:"risk_owner_first_name_en"`
	RiskOwnerFirstNameAr          string     `json:"risk_owner_first_name_ar"`
	RiskOwnerLastNameEn           string     `json:"risk_owner_last_name_en"`
	RiskOwnerLastNameAr           string     `json:"risk_owner_last_name_ar"`
	RiskOwnerInitial              string     `json:"risk_owner_initial"`
	RiskOwnerColor                string     `json:"risk_owner_color"`
	RiskOwnerEmail                string     `json:"risk_owner_email"`
	RiskOwnerProfileImagePath     string     `json:"risk_owner_profile_image_path"`
	RiskReportedAt                *time.Time `json:"risk_reported_at"`
	RiskReviewFrequencyTitleEn    *string    `json:"risk_review_frequency_title_en"`
	RiskReviewFrequencyTitleAr    *string    `json:"risk_review_frequency_title_ar"`
	RiskNextReviewDate            *string    `json:"risk_next_review_date"`
	RiskStatusTitleEn             string     `json:"risk_status_title_en"`
	RiskStatusTitleAr             string     `json:"risk_status_title_ar"`
	RiskStatusColorCode           string     `json:"risk_status_color_code"`
	RiskMatrixLikelihoodTitleEn   string     `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAr   string     `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixLikelihoodColorCode string     `json:"risk_matrix_likelihood_score"`
	RiskMatrixImpactTitleEn       string     `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr       string     `json:"risk_matrix_impact_title_ar"`
	RiskMatrixImpactColorCode     string     `json:"risk_matrix_impact_score"`
	RiskMatrixRiskRatingTitleEn   string     `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr   string     `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingType      string     `json:"risk_matrix_risk_rating_type"`
	RiskMatrixRiskRatingColorCode string     `json:"risk_matrix_risk_rating_color_code"`
}

type ComplianceRequirementIncidentListDTO struct {
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
	IncidentReportedUserFirstNameEn      string `json:"incident_reported_user_first_name_en"`
	IncidentReportedUserFirstNameAr      string `json:"incident_reported_user_first_name_ar"`
	IncidentReportedUserLastNameEn       string `json:"incident_reported_user_last_name_en"`
	IncidentReportedUserLastNameAr       string `json:"incident_reported_user_last_name_ar"`
	IncidentReportedUserInitial          string `json:"incident_reported_user_initial"`
	IncidentReportedUserColor            string `json:"incident_reported_user_color"`
	IncidentReportedUserEmail            string `json:"incident_reported_user_email"`
	IncidentReportedUserProfileImagePath string `json:"incident_reported_user_profile_image_path"`
}

type ComplianceRequirementAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type ComplianceRequirementAssociatedItemControlListResponseDTO struct {
	Items      []*ComplianceRequirementControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ComplianceRequirementAssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*ComplianceRequirementStrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO                      `json:"pagination"`
}

type ComplianceRequirementAssociatedItemRiskListResponseDTO struct {
	Items      []*ComplianceRequirementRiskListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

type ComplianceRequirementAssociatedItemIncidentListResponseDTO struct {
	Items      []*ComplianceRequirementIncidentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO            `json:"pagination"`
}
