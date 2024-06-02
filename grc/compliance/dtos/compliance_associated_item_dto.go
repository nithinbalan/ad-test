package dtos

type ComplianceControlListDTO struct {
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

type ComplianceAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}
