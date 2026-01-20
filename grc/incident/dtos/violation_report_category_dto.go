package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
)

type ViolationReportCategoryFormDTO struct {
	DepartmentID  uint    `json:"department_id" binding:"required"`
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	Status        string  `json:"status" binding:"required,oneof=active inactive"`
}

type ViolationReportCategoryListDTO struct {
	ID                uint    `json:"id"`
	DepartmentID      uint    `json:"department_id"`
	DepartmentTitleEn string  `json:"department_title_en"`
	DepartmentTitleAr string  `json:"department_title_ar"`
	TitleEn           string  `json:"title_en"`
	TitleAr           string  `json:"title_ar"`
	DescriptionEn     *string `json:"description_en,omitempty"`
	DescriptionAr     *string `json:"description_ar,omitempty"`
	Status            string  `json:"status"`
}

type ViolationReportCategoryListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type ViolationReportCategoryDetailsDTO struct {
	ID            uint                                  `json:"id"`
	Department    *organizationDTOs.OrganizationUnitDTO `json:"department"`
	TitleEn       string                                `json:"title_en"`
	TitleAr       string                                `json:"title_ar"`
	DescriptionEn string                                `json:"description_en,omitempty"`
	DescriptionAr string                                `json:"description_ar,omitempty"`
	Status        string                                `json:"status"`
}
