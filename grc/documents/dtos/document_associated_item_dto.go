package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentComplianceRequirementListDTO struct {
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

type DocumentProcessListDTO struct {
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

type DocumentStrategicObjectiveListDTO struct {
	StrategicObjectiveID      uint    `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr *string `json:"strategic_objective_title_ar"`
	StrategicObjectiveStatus  *string `json:"strategic_objective_status"`
}
type DocumentSLAListDTO struct {
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
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
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

type DocumentIncidentListDTO struct {
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

type DocumentRiskListDTO struct {
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr             string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn               string  `json:"risk_source_title_en"`
	RiskSourceTitleAr               string  `json:"risk_source_title_ar"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ResponsibleUserUnitTitleEn      string  `json:"responsible_user_unit_title_en"`
	ResponsibleUserUnitTitleAr      string  `json:"responsible_user_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	RiskMatrixLikelihoodTitleEn     string  `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAr     string  `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixLikelihoodColorCode   string  `json:"risk_matrix_likelihood_score"`
	RiskMatrixImpactTitleEn         string  `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr         string  `json:"risk_matrix_impact_title_ar"`
	RiskMatrixImpactColorCode       string  `json:"risk_matrix_impact_score"`
	RiskMatrixRiskRatingTitleEn     string  `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr     string  `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingType        string  `json:"risk_matrix_risk_rating_type"`
	RiskMatrixRiskRatingColorCode   string  `json:"risk_matrix_risk_rating_color_code"`
}

type DocumentAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type DocumentAssociatedItemControlListResponseDTO struct {
	Items      []*commonDTOs.AssociatedControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type DocumentAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*organizationDTOs.ComplianceRegisterDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO              `json:"pagination"`
}

type DocumentAssociatedItemProcessListResponseDTO struct {
	Items      []*DocumentProcessListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type DocumentAssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*DocumentStrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO         `json:"pagination"`
}

type DocumentAssociatedItemSLAListResponseDTO struct {
	Items      []*DocumentSLAListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type DocumentAssociatedItemIncidentListResponseDTO struct {
	Items      []*DocumentIncidentListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type DocumentAssociatedItemRiskListResponseDTO struct {
	Items      []*DocumentRiskListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ClauseDTO struct {
	ComplianceClauseID uint   `json:"compliance_clause_id"`
	ClauseNo           string `json:"clause_no"`
	ClauseTitleEn      string `json:"clause_title_en"`
	ClauseTitleAr      string `json:"clause_title_ar"`
}

type ClauseListDTO struct {
	ComplianceRegisterID      uint        `json:"compliance_register_id"`
	ComplianceRegisterTitleEn string      `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr string      `json:"compliance_register_title_ar"`
	Clauses                   []ClauseDTO `json:"clauses"`
}
type ClauseRowDTO struct {
	ComplianceRegisterID      uint
	ComplianceRegisterTitleEn string
	ComplianceRegisterTitleAr string
	ComplianceClauseID        uint
	ClauseNo                  string
	ClauseTitleEn             string
	ClauseTitleAr             string
}
type DocumentAssociatedItemClauseListResponseDTO struct {
	Items      []*ClauseListDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
type DocumentAssociatedItemClauseFormDTO struct {
	ClauseIDs []uint `json:"clause_ids" binding:"required"`
}
