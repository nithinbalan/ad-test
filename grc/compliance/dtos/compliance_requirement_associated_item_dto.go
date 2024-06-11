package dtos

type ComplianceRequirementRiskListDTO struct {
	RiskRegisterID      uint   `json:"risk_register_id"`
	RiskRegisterTitle   string `json:"risk_register_title"`
	RiskStatusTitleEn   string `json:"risk_status_title_en"`
	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
	RiskStatusColorCode string `json:"risk_status_color_code"`
}

type ComplianceRequirementIncidentListDTO struct {
	IncidentRegisterID          uint   `json:"incident_register_id"`
	IncidentRegisterTitle       string `json:"incident_register_title"`
	IncidentRegisterDescription string `json:"incident_register_description"`
	IncidentStatusTitleEn       string `json:"incident_status_title_en"`
	IncidentStatusTitleAr       string `json:"incident_status_title_ar"`
	IncidentStatusColorCode     string `json:"incident_status_color_code"`
}

type ComplianceRequirementAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}
