package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceRequirementListDTO struct {
	ID                                   uint   `json:"id"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	Title                                string `json:"title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	SeverityLevelTitleEn                 string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                 string `json:"severity_level_title_ar"`
	SeverityLevelColorCode               string `json:"severity_level_color_code"`
	ResponsibleUserFirstName             string `json:"responsible_user_first_name"`
	ResponsibleUserLastName              string `json:"responsible_user_last_name"`
	Deadline                             string `json:"deadline"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
	CreatedAt                            string `json:"created_at"`
	UpdatedAt                            string `json:"updated_at"`
	CreatedUserFirstNameEn               string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn               string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                string `json:"updated_user_last_name_ar"`
}

type ComplianceRequirementListQueryParamsDTO struct {
	Sort         string
	SortBy       string
	Keyword      string
	Status       string
	Page         int
	PageSize     int
	Filters      map[string][]string
	QuickFilter  []string
	DeadlineFrom string
	DeadlineTo   string
}

type ComplianceRequirementListResponseDTO struct {
	Items      []*ComplianceRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type ComplianceRequirementFormDTO struct {
	ComplianceClauseNo              string `json:"compliance_clause_no"`
	ComplianceClauseTitle           string `json:"compliance_clause_title"`
	Title                           string `json:"title"`
	ComplianceRequirementCategoryID uint   `json:"compliance_requirement_category_id" binding:"required"`
	SeverityLevelID                 uint   `json:"severity_level_id"`
	ResponsibleUserID               uint   `json:"responsible_user_id" binding:"required"`
	Deadline                        string `json:"deadline"`
	ComplianceStatusID              uint   `json:"compliance_status_id" binding:"required"`
}
