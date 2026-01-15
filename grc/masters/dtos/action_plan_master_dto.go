package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type ActionPlanMasterListDTO struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      string  `json:"status"`

	ModuleID         uint   `json:"module_id"`
	ModuleIdentifier string `json:"module_identifier"`
	ModuleTitleEn    string `json:"module_title_en"`
	ModuleTitleAr    string `json:"module_title_ar"`

	IncidentCategoryID      *uint   `json:"incident_category_id"`
	IncidentCategoryTitleEn *string `json:"incident_category_title_en"`
	IncidentCategoryTitleAr *string `json:"incident_category_title_ar"`

	SeverityLevelID   *uint   `json:"severity_level_id"`
	SeverityTitleEn   *string `json:"severity_title_en"`
	SeverityTitleAr   *string `json:"severity_title_ar"`
	SeverityColorCode *string `json:"severity_color_code"`
	SeverityType      *string `json:"severity_type"`
}

type ActionPlanMasterDetailsDTO struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      string  `json:"status"`

	Module           ModuleDTO            `json:"module"`
	IncidentCategory *IncidentCategoryDTO `json:"incident_category"`
	SeverityLevel    *SeverityLevelDTO    `json:"severity_level"`
}

type ModuleDTO struct {
	ID         uint   `json:"id"`
	Identifier string `json:"identifier" `
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
}

type IncidentCategoryDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type ActionPlanMasterFormDTO struct {
	ModuleID           uint    `json:"module_id" binding:"required"`
	IncidentCategoryID *uint   `json:"incident_category_id"`
	SeverityLevelID    *uint   `json:"severity_level_id"`
	Title              string  `json:"title" binding:"required"`
	Description        *string `json:"description"`
}

type ActionPlanMasterListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Status   string
	Page     int
	PageSize int
}

type ActionPlanMasterListResponseDTO struct {
	Items      []*ActionPlanMasterListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
