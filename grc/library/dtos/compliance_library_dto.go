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

type ClauseWithSubClausesDTO struct {
	ID            uint                       `json:"id"`
	ClauseNumber  string                     `json:"clause_number"`
	ClauseTitleEn string                     `json:"clause_title_en"`
	ClauseTitleAr string                     `json:"clause_title_ar"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedAt     string                     `json:"updated_at"`
	SubClauses    []SubClauseWithControlsDTO `json:"sub_clauses"`
}

type SubClauseWithControlsDTO struct {
	ID                  uint         `json:"id"`
	ClauseNumber        string       `json:"clause_number"`
	ClauseTitleEn       string       `json:"clause_title_en"`
	ClauseTitleAr       string       `json:"clause_title_ar"`
	ClauseDescriptionEn string       `json:"clause_description_en"`
	ClauseDescriptionAr string       `json:"clause_description_ar"`
	Controls            []ControlDTO `json:"controls"`
}

type ControlDTO struct {
	ID             uint   `json:"id"`
	ControlCode    string `json:"control_code"`
	ControlTitleEn string `json:"control_title_en"`
	ControlTitleAr string `json:"control_title_ar"`
}
