package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingControlListDTO struct {
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

type MeetingComplianceRequirementListDTO struct {
	ComplianceRegisterID                 uint   `json:"compliance_register_id"`
	ComplianceRegisterTitleEn            string `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr            string `json:"compliance_register_title_ar"`
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitleEn              string `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr              string `json:"compliance_clause_title_ar"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitleEn         string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string `json:"compliance_requirement_title_ar"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type MeetingProcessListDTO struct {
	ProcessGroupTitleEn      string  `json:"process_group_title_en"`
	ProcessGroupTitleAr      string  `json:"process_group_title_ar"`
	OwnerUnitTitleEn         *string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr         *string `json:"owner_unit_title_ar"`
	OwnerID                  *uint   `json:"owner_id"`
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

type MeetingRiskListDTO struct {
	ReferenceNumber               string  `json:"reference_number"`
	RiskCategoryTitleEn           string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr           string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn             string  `json:"risk_source_title_en"`
	RiskSourceTitleAr             string  `json:"risk_source_title_ar"`
	RiskStatusTitleEn             string  `json:"risk_status_title_en"`
	RiskStatusTitleAr             string  `json:"risk_status_title_ar"`
	RiskStatusColorCode           string  `json:"risk_status_color_code"`
	Title                         string  `json:"title"`
	Description                   *string `json:"description"`
	ResponsibleUnitTitleEn        string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr        string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID             uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn    string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr    string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn     string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr     string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial        string  `json:"responsible_user_initial"`
	ResponsibleUserColor          string  `json:"responsible_user_color"`
	ResponsibleUserEmail          string  `json:"responsible_user_email"`
	ResponsibleProfileImagePath   string  `json:"responsible_user_profile_image_path"`
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

type MeetingIncidentListDTO struct {
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

type MeetingAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type MeetingAssociatedItemControlListResponseDTO struct {
	Items      []*MeetingControlListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*organizationDTOs.ComplianceRegisterDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO              `json:"pagination"`
}

type MeetingAssociatedItemProcessListResponseDTO struct {
	Items      []*MeetingProcessListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingAssociatedItemRiskListResponseDTO struct {
	Items      []*MeetingRiskListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingAssociatedItemIncidentListResponseDTO struct {
	Items      []*MeetingIncidentListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
