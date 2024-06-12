package dtos

import "time"

type ProcessControlListDTO struct {
	ControlID                       uint    `json:"control_id"`
	ControlCategoryTitleEn          string  `json:"control_category_title_en"`
	ControlCategoryTitleAr          string  `json:"control_category_title_ar"`
	ControlTypeTitleEn              string  `json:"control_type_title_en"`
	ControlTypeTitleAr              string  `json:"control_type_title_ar"`
	ControlTitleEn                  string  `json:"control_title_en"`
	ControlTitleAr                  string  `json:"control_title_ar"`
	ControlDescription              *string `json:"control_description"`
	ControlEfficiencyStatusEn       string  `json:"control_efficiency_status_en"`
	ControlEfficiencyStatusAr       string  `json:"control_efficiency_status_ar"`
	ControlEfficiencyStatusType     string  `json:"control_efficiency_status_type"`
	ControlImplementationStatusEn   string  `json:"control_implementation_status_en"`
	ControlImplementationStatusAr   string  `json:"control_implementation_status_ar"`
	ControlImplementationStatusType string  `json:"control_implementation_status_type"`
}

type ProcessRiskListDTO struct {
	ReferenceNumber     string  `json:"reference_number"`
	RiskCategoryTitleEn string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn   string  `json:"risk_source_title_en"`
	RiskSourceTitleAr   string  `json:"risk_source_title_ar"`
	OwnerUnitTitleEn    string  `json:"owner_unit_title_en"`
	OwnerUnitTitleAr    string  `json:"owner_unit_title_ar"`
	OwnerFirstNameEn    string  `json:"owner_first_name_en"`
	OwnerFirstNameAr    string  `json:"owner_first_name_ar"`
	OwnerLastNameEn     string  `json:"owner_last_name_en"`
	OwnerLastNameAr     string  `json:"owner_last_name_ar"`
	RiskStatusTitleEn   string  `json:"risk_status_title_en"`
	RiskStatusTitleAr   string  `json:"risk_status_title_ar"`
	RiskStatusColorCode string  `json:"risk_status_color_code"`
	Title               string  `json:"title"`
	Description         *string `json:"description"`
}

type ProcessIncidentListDTO struct {
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
	QuickFilter []string
}
