package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type UserListDTO struct {
	ID                       uint   `json:"id"`
	FirstNameEn              string `json:"first_name_en"`
	LastNameEn               string `json:"last_name_en"`
	FirstNameAr              string `json:"first_name_ar"`
	LastNameAr               string `json:"last_name_ar"`
	OrganizationUnit1TitleEn string `json:"organization_unit_1_title_en"`
	OrganizationUnit2TitleEn string `json:"organization_unit_2_title_en,omitempty"`
	OrganizationUnit3TitleEn string `json:"organization_unit_3_title_en,omitempty"`
	OrganizationUnit4TitleEn string `json:"organization_unit_4_title_en,omitempty"`
	OrganizationUnit5TitleEn string `json:"organization_unit_5_title_en,omitempty"`
	OrganizationUnit1TitleAr string `json:"organization_unit_1_title_ar"`
	OrganizationUnit2TitleAr string `json:"organization_unit_2_title_ar,omitempty"`
	OrganizationUnit3TitleAr string `json:"organization_unit_3_title_ar,omitempty"`
	OrganizationUnit4TitleAr string `json:"organization_unit_4_title_ar,omitempty"`
	OrganizationUnit5TitleAr string `json:"organization_unit_5_title_ar,omitempty"`
	DesignationTitleEn       string `json:"designation_title_en"`
	DesignationTitleAr       string `json:"designation_title_ar"`
}

type UserListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type UserListResponseDTO struct {
	Items      []*UserListDTO               `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type UserFormDTO struct {
	OrganizationUnit1ID uint    `json:"organization_unit_1_id" binding:"required"`
	OrganizationUnit2ID *uint   `json:"organization_unit_2_id" binding:"omitempty"`
	OrganizationUnit3ID *uint   `json:"organization_unit_3_id" binding:"omitempty"`
	OrganizationUnit4ID *uint   `json:"organization_unit_4_id" binding:"omitempty"`
	OrganizationUnit5ID *uint   `json:"organization_unit_5_id" binding:"omitempty"`
	DesignationID       uint    `json:"designation_id" binding:"required"`
	FirstNameEn         string  `json:"first_name_en" binding:"required,max=255"`
	LastNameEn          string  `json:"last_name_en" binding:"required,max=255"`
	FirstNameAr         string  `json:"first_name_ar" binding:"required,max=255"`
	LastNameAr          string  `json:"last_name_ar" binding:"required,max=255"`
	Email               string  `json:"email" binding:"required,email"`
	MobileNumber        *string `json:"mobile_number" binding:"omitempty,max=45"`
	ProfileImagePath    *string `json:"profile_image_path" binding:"omitempty"`
	RoleIDs             []uint  `json:"role_ids" binding:"required"`
}
