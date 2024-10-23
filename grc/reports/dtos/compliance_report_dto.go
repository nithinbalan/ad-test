package dtos

type ComplianceByCategoryDTO struct {
	ComplianceCategoryID      uint   `json:"compliance_category_id"`
	ComplianceCategoryTitleEn string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr string `json:"compliance_category_title_ar"`
	Count                     int64  `json:"count"`
}

type ReportListQueryParamsDTO struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
