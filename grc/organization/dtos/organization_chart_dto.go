package dtos

import (
	settingDTOs "jethur/grc/settings/dtos"
)

type OrganizationChartUserListDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameAr       string  `json:"last_name_ar"`
	Email            string  `json:"email"`
	MobileNumber     *string `json:"mobile_number"`
	ProfileImagePath *string `json:"profile_image_path"`
	Initial          string  `json:"initial"`
	Color            string  `json:"color"`

	DesignationEn             string `json:"designation_en"`
	DesignationAr             string `json:"designation_ar"`
	PrimaryOrganizationUnitID *uint  `json:"primary_organization_unit_id" gorm:"column:primary_organization_unit_id"`
	PrimaryOrganizationUnitEn string `json:"primary_organization_unit_en" gorm:"column:primary_organization_unit_en"`
	PrimaryOrganizationUnitAr string `json:"primary_organization_unit_ar" gorm:"column:primary_organization_unit_ar"`

	ReportingUserID          *uint  `json:"reporting_user_id"`
	ReportingUserFirstNameEn string `json:"reporting_user_first_name_en"`
	ReportingUserLastNameEn  string `json:"reporting_user_last_name_en"`
	ReportingUserFirstNameAr string `json:"reporting_user_first_name_ar"`
	ReportingUserLastNameAr  string `json:"reporting_user_last_name_ar"`

	RolesEn             string `json:"roles_en"`
	RolesAr             string `json:"roles_ar"`
	ReportingUsersCount int    `json:"reporting_users_count"`
}

type OrganizationChartUserDetailsDTO struct {
	ID                        uint            `json:"id"`
	FirstNameEn               string          `json:"first_name_en"`
	LastNameEn                string          `json:"last_name_en"`
	FirstNameAr               string          `json:"first_name_ar"`
	LastNameAr                string          `json:"last_name_ar"`
	Email                     string          `json:"email"`
	MobileNumber              *string         `json:"mobile_number"`
	ProfileImagePath          *string         `json:"profile_image_path"`
	IsHead                    bool            `json:"is_head"`
	Designation               *DesignationDTO `json:"designation"`
	PrimaryOrganizationUnitID *uint           `json:"primary_organization_unit_id"`
	PrimaryOrganizationUnitEn string          `json:"primary_organization_unit_en"`
	PrimaryOrganizationUnitAr string          `json:"primary_organization_unit_ar"`

	Roles         *[]settingDTOs.RoleDTO `json:"roles"`
	ReportingUser *UserMinimalDTO        `json:"reporting_user"`
}

type OrganizationChartUserListQueryParamsDTO struct {
	Filters map[string][]string
}

type ChangeReportingUserFormDTO struct {
	ReportingUserID *uint `json:"reporting_user_id"`
}
