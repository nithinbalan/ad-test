package dtos

type ViolationReportSubCategoryFormDTO struct {
	ViolationReportCategoryID uint    `json:"violation_report_category_id" binding:"required"`
	TitleEn                   string  `json:"title_en" binding:"required"`
	TitleAr                   string  `json:"title_ar" binding:"required"`
	DescriptionEn             *string `json:"description_en"`
	DescriptionAr             *string `json:"description_ar"`
	Status                    string  `json:"status" binding:"required,oneof=active inactive"`
}

type ViolationReportSubCategoryListDTO struct {
	ID                             uint    `json:"id"`
	ViolationReportCategoryID      uint    `json:"violation_report_category_id"`
	ViolationReportCategoryTitleEn string  `json:"violation_report_category_title_en"`
	ViolationReportCategoryTitleAr string  `json:"violation_report_category_title_ar"`
	TitleEn                        string  `json:"title_en"`
	TitleAr                        string  `json:"title_ar"`
	DescriptionEn                  *string `json:"description_en,omitempty"`
	DescriptionAr                  *string `json:"description_ar,omitempty"`
	Status                         string  `json:"status"`
}

type ViolationReportSubCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type ViolationReportSubCategoryDetailsDTO struct {
	ID                      uint                              `json:"id"`
	ViolationReportCategory ViolationReportCategoryDetailsDTO `json:"violation_report_category"`
	TitleEn                 string                            `json:"title_en"`
	TitleAr                 string                            `json:"title_ar"`
	DescriptionEn           string                            `json:"description_en,omitempty"`
	DescriptionAr           string                            `json:"description_ar,omitempty"`
	Status                  string                            `json:"status"`
}
