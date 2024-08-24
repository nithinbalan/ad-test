package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type SLAProcessListDTO struct {
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

type SLADocumentListDTO struct {
	ID                         uint       `json:"id"`
	ReferenceNumber            string     `json:"reference_number"`
	TitleEn                    string     `json:"title_en"`
	TitleAr                    string     `json:"title_ar"`
	DescriptionEn              *string    `json:"description_en"`
	DescriptionAr              *string    `json:"description_ar"`
	DocumentCategoryTitleEn    string     `json:"document_category_title_en"`
	DocumentCategoryTitleAr    string     `json:"document_category_title_ar"`
	ResponsibleUserFirstNameEn string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn  string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr  string     `json:"responsible_user_last_name_ar"`
	ResponsibleUnitTitleEn     string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr     string     `json:"responsible_unit_title_ar"`
	DocumentStatusTitleEn      string     `json:"document_status_title_en"`
	DocumentStatusTitleAr      string     `json:"document_status_title_ar"`
	DocumentStatusColorCode    string     `json:"document_status_color_code"`
	LastCheckoutAt             *time.Time `json:"last_checkout_at"`
	LastCheckinAt              *time.Time `json:"last_checkin_at"`
	IsLocked                   bool       `json:"is_locked"`
}

type SLAIncidentListDTO struct {
	ID                         uint    `json:"id"`
	IncidentRegisterID         uint    `json:"incident_register_id"`
	IncidentReferenceNumber    string  `json:"incident_reference_number"`
	IncidentTitle              string  `json:"incident_title"`
	IncidentDescription        *string `json:"incident_description"`
	IncidentCategoryTitleEn    string  `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string  `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn *string `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr *string `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn        string  `json:"incident_type_title_en"`
	IncidentTypeTitleAr        string  `json:"incident_type_title_ar"`
	SeverityLevelTitleEn       string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode     string  `json:"severity_level_color_code"`
	IncidentStatusTitleEn      string  `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string  `json:"incident_status_title_ar"`
	IncidentStatusColorCode    string  `json:"incident_status_color_code"`
}

type SLARiskListDTO struct {
	RiskRegisterID                  uint    `json:"risk_register_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr             string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string `json:"risk_source_title_ar"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
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

type SLAAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type SLAAssociatedItemProcessListResponseDTO struct {
	Items      []*SLAProcessListDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLAAssociatedItemDocumentListResponseDTO struct {
	Items      []*SLADocumentListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLAAssociatedItemIncidentListResponseDTO struct {
	Items      []*SLAIncidentListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLAAssociatedItemRiskListResponseDTO struct {
	Items      []*SLARiskListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SlaAssociatedItemsCountResponseDTO struct {
	Questions   int `json:"questions"`
	Risks       int `json:"risks"`
	Incidents   int `json:"incidents"`
	ActionPlans int `json:"action_plans"`
}
