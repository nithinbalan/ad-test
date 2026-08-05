package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
)

type RoleDTO struct {
	ID                   uint    `json:"id"`
	TitleEn              string  `json:"title_en"`
	TitleAr              string  `json:"title_ar"`
	Status               string  `json:"status"`
	IsSuperAdmin         bool    `json:"is_super_admin"`
	IsCyber              bool    `json:"is_cyber"`
	LicenseTypeID        uint    `json:"license_type_id"`
	LicenseTypeTitleEn   string  `json:"license_type_title_en"`
	LicenseTypeTitleAr   string  `json:"license_type_title_ar"`
	LicenseTypeType      string  `json:"license_type_type"`
	LicenseTypeColorCode string  `json:"license_type_color_code"`
	Type                 *string `json:"type"`
}

type RoleDetailDTO struct {
	ID            uint                      `json:"id"`
	TitleEn       string                    `json:"title_en"`
	TitleAr       string                    `json:"title_ar"`
	Status        string                    `json:"status"`
	IsSuperAdmin  bool                      `json:"is_super_admin"`
	IsCyber       bool                      `json:"is_cyber"`
	LicenseTypeID uint                      `json:"license_type_id"`
	LicenseType   masterDTOs.LicenseTypeDTO `json:"license_type"`
	Type          *string                   `json:"type"`
}

type RolePermissionFormDTO struct {
	ModuleID             uint `json:"module_id" binding:"required"`
	RolePermissionTypeID uint `json:"role_permission_type_id" binding:"required"`
	IsGranted            bool `json:"is_granted" binding:"required"`
}
