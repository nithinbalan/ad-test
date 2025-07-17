package dtos

import commonDTOs "jethur/grc/common/dtos"

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

type ControlItemDTO struct {
	ControlCode    string `json:"control_code"`
	ControlTitleEn string `json:"control_title_en"`
	ControlTitleAr string `json:"control_title_ar"`
}

type FrameworkControlsDTO struct {
	Framework string           `json:"framework"`
	Controls  []ControlItemDTO `json:"controls"`
}

type ControlMappingDTO struct {
	ID         string                 `json:"id"` // example: CTRL-001
	TitleEn    string                 `json:"title_en"`
	TitleAr    string                 `json:"title_ar"`
	Frameworks []FrameworkControlsDTO `json:"frameworks"`
}

type RawControlMapping struct {
	GenericControlID      uint
	GenericControlTitleEn string
	GenericControlTitleAr string
	ControlCode           string
	ControlTitleEn        string
	ControlTitleAr        string
	FrameworkEn           string
	FrameworkAr           string
}

type ControlLibraryResponseDTO struct {
	ID      uint              `json:"id"`
	TitleEn string            `json:"title_en"`
	TitleAr string            `json:"title_ar"`
	SAMA    []ControlTitleDTO `json:"sama,omitempty" gorm:"-"`
	NCA     []ControlTitleDTO `json:"nca,omitempty" gorm:"-"`
}

type ControlTitleDTO struct {
	ControlLibraryID         uint   `json:"control_library_id"`
	ComplianceLibraryTitleEn string `json:"compliance_library_title_en"`
	ComplianceLibraryTitleAr string `json:"compliance_library_title_ar"`
	ControlCode              string `json:"control_code"`
	ControlTitleEn           string `json:"control_title_en"`
	ControlTitleAr           string `json:"control_title_ar"`
}

type ControlMappingListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}

type ControlMappingListResponseDTO struct {
	Items      []*ControlLibraryResponseDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
