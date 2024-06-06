package dtos

type IncidentControlListDTO struct {
	ControlID                     uint    `json:"control_id"`
	ControlCategoryTitleEn        string  `json:"control_category_title_en"`
	ControlCategoryTitleAr        string  `json:"control_category_title_ar"`
	ControlTypeTitleEn            string  `json:"control_type_title_en"`
	ControlTypeTitleAr            string  `json:"control_type_title_ar"`
	ControlTitleEn                string  `json:"control_title_en"`
	ControlTitleAr                string  `json:"control_title_ar"`
	ControlDescription            *string `json:"control_description"`
	ControlEfficiencyStatusEn     string  `json:"control_efficiency_status_en"`
	ControlEfficiencyStatusAr     string  `json:"control_efficiency_status_ar"`
	ControlImplementationStatusEn string  `json:"control_implementation_status_en"`
	ControlImplementationStatusAr string  `json:"control_implementation_status_ar"`
}

type IncidentComplianceRequirementListDTO struct {
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitle           string `json:"compliance_requirement_title"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type IncidentRiskListDTO struct {
	RiskRegisterID      uint   `json:"risk_register_id"`
	RiskRegisterTitle   string `json:"risk_register_title"`
	RiskStatusTitleEn   string `json:"risk_status_title_en"`
	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
	RiskStatusColorCode string `json:"risk_status_color_code"`
}

type IncidentAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}
