package dtos

import "time"

type ComplianceRequirementRiskListDTO struct {
	RiskRegisterID             uint       `json:"risk_register_id"`
	RiskReferenceNumber        string     `json:"risk_reference_number"`
	RiskRegisterTitle          string     `json:"risk_register_title"`
	RiskRegisterDescription    string     `json:"risk_register_description"`
	RiskCategoryTitleEn        string     `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string     `json:"risk_category_title_ar"`
	RiskSourceTitleEn          string     `json:"risk_source_title_en"`
	RiskSourceTitleAr          string     `json:"risk_source_title_ar"`
	RiskOwnerUnitTitleEn       string     `json:"risk_owner_unit_title_en"`
	RiskOwnerUnitTitleAr       string     `json:"risk_owner_unit_title_ar"`
	RiskOwnerFirstNameEn       string     `json:"risk_owner_first_name_en"`
	RiskOwnerFirstNameAr       string     `json:"risk_owner_first_name_ar"`
	RiskOwnerLastNameEn        string     `json:"risk_owner_last_name_en"`
	RiskOwnerLastNameAr        string     `json:"risk_owner_last_name_ar"`
	RiskReportedAt             *time.Time `json:"risk_reported_at"`
	RiskReviewFrequencyTitleEn *string    `json:"risk_review_frequency_title_en"`
	RiskReviewFrequencyTitleAr *string    `json:"risk_review_frequency_title_ar"`
	RiskNextReviewDate         *string    `json:"risk_next_review_date"`
	RiskStatusTitleEn          string     `json:"risk_status_title_en"`
	RiskStatusTitleAr          string     `json:"risk_status_title_ar"`
	RiskStatusColorCode        string     `json:"risk_status_color_code"`
}

type ComplianceRequirementIncidentListDTO struct {
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

type ComplianceRequirementAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}
