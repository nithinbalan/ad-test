package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type TenableAssetListQueryParamsDTO struct {
	Page     int
	PageSize int
}

type TenableAssetDTO struct {
	TenableUUID string `json:"tenable_uuid" binding:"required"`
	TitleEn     string `json:"title_en" binding:"required"`
	TitleAr     string `json:"title_ar" binding:"required"`

	AssetCategoryID uint    `json:"asset_category_id" binding:"required"`
	AssetCategory   *string `json:"asset_category,omitempty"`

	AssetSubCategoryID uint    `json:"asset_sub_category_id" binding:"required"`
	AssetSubCategory   *string `json:"asset_sub_category,omitempty"`

	AssetStatusID uint    `json:"asset_status_id" binding:"required"`
	AssetStatus   *string `json:"asset_status,omitempty"`

	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	ResponsibleUser   *string `json:"responsible_user,omitempty"`

	Location           *string                       `json:"location" binding:"omitempty"`
	Description        *string                       `json:"description" binding:"omitempty"`
	AcquisitionDate    *string                       `json:"acquisition_date" binding:"omitempty,datetime=2006-01-02"`
	AssignmentDate     *string                       `json:"assignment_date" binding:"omitempty,datetime=2006-01-02"`
	WarrantyExpiryDate *string                       `json:"warranty_expiry_date" binding:"omitempty,datetime=2006-01-02"`
	Manufacturer       *string                       `json:"manufacturer" binding:"omitempty"`
	Model              *string                       `json:"model" binding:"omitempty"`
	SerialNumber       *string                       `json:"serial_number" binding:"omitempty"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents"`
}

type TenableAssetListResponseDTO struct {
	Items      []TenableAssetDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
