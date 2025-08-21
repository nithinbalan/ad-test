package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ControlIncidentListDTO struct {
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

type ControlAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}

type ControlAssociatedItemIncidentListResponseDTO struct {
	Items      []*ControlIncidentListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlComplianceRequirementListDTO struct {
	ComplianceRegisterID                 uint   `json:"compliance_register_id"`
	ComplianceRegisterTitleEn            string `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr            string `json:"compliance_register_title_ar"`
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitleEn              string `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr              string `json:"compliance_clause_title_ar"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitleEn         string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string `json:"compliance_requirement_title_ar"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type ComplianceRequirementMinimalDTO struct {
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitleEn              string `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr              string `json:"compliance_clause_title_ar"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitleEn         string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string `json:"compliance_requirement_title_ar"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type ComplianceRegisterDTO struct {
	ComplianceRegisterID      uint                               `json:"compliance_register_id"`
	ComplianceRegisterTitleEn string                             `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr string                             `json:"compliance_register_title_ar"`
	Requirements              []*ComplianceRequirementMinimalDTO `json:"requirements"`
}

type ControlAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*ComplianceRegisterDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlAssociatedItemsCountResponseDTO struct {
	Processes              int `json:"processes"`
	ComplianceRequirements int `json:"compliance_requirements"`
	Risks                  int `json:"risks"`
	Incidents              int `json:"incidents"`
}

type ClauseDTO struct {
	ComplianceClauseID uint   `json:"compliance_clause_id"`
	ClauseNo           string `json:"clause_no"`
	ClauseTitleEn      string `json:"clause_title_en"`
	ClauseTitleAr      string `json:"clause_title_ar"`
}

type ClauseListDTO struct {
	ComplianceRegisterID      uint        `json:"compliance_register_id"`
	ComplianceRegisterTitleEn string      `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr string      `json:"compliance_register_title_ar"`
	Clauses                   []ClauseDTO `json:"clauses"`
}
type ClauseRowDTO struct {
	ComplianceRegisterID      uint
	ComplianceRegisterTitleEn string
	ComplianceRegisterTitleAr string
	ComplianceClauseID        uint
	ClauseNo                  string
	ClauseTitleEn             string
	ClauseTitleAr             string
}
type ControlAssociatedItemClauseListResponseDTO struct {
	Items      []*ClauseListDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
type ControlClauseFormDTO struct {
	ClauseIDs []uint `json:"clause_ids" binding:"required"`
}
