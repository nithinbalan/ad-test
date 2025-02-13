package dtos

type ComplianceFolderListDTO struct {
	ID                      uint   `json:"id"`
	TitleEn                 string `json:"title_en"`
	TitleAr                 string `json:"title_ar"`
	Status                  string `json:"status"`
	ComplianceRegisterCount int    `json:"compliance_register_count"`
}

type ComplianceFolderDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
