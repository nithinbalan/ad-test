package dtos

type ComplianceLibraryWithVersionsDTO struct {
	ID                     uint    `json:"id"`
	Icon                   *string `json:"icon"`
	TitleEn                string  `json:"title_en"`
	TitleAr                string  `json:"title_ar"`
	DescriptionEn          *string `json:"description_en"`
	DescriptionAr          *string `json:"description_ar"`
	IssuedBy               *string `json:"issued_by"`
	Classification         *string `json:"classification"`
	DocumentClassification *string `json:"document_classification"`
	Disclaimer             *string `json:"disclaimer"`

	VersionID    *uint   `json:"version_id"`
	Version      *string `json:"version"`
	Year         *int    `json:"year"`
	Overview     *string `json:"overview"`
	Introduction *string `json:"introduction"`
	Domains      *string `json:"domains"`
}
