package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AssetCriticalityMatrixListDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type AssetCriticalityMatrixDetailDTO struct {
	ID              uint                                          `json:"id"`
	TitleEn         string                                        `json:"title_en"`
	TitleAr         string                                        `json:"title_ar"`
	Status          string                                        `json:"status"`
	AssetCategories []*AssetCriticalityMatrixAssetCategoryListDTO `json:"asset_categories"`
}

type AssetCriticalityMatrixFormDTO struct {
	TitleEn          string `json:"title_en" binding:"required"`
	TitleAr          string `json:"title_ar" binding:"required"`
	Status           string `json:"status" binding:"required"`
	AssetCategoryIDs []uint `json:"asset_category_ids" binding:"required"`
}

type AssetCriticalityMatrixListQueryParamsDTO struct {
	Sort     string `json:"sort,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Status   string `json:"status,omitempty"`
}

type AssetCriticalityMatrixListResponseDTO struct {
	Items      []*AssetCriticalityMatrixListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type AssetCriticalityMatrixAssetCategoryListDTO struct {
	ID                       uint             `json:"id"`
	AssetCriticalityMatrixID uint             `json:"asset_criticality_matrix_id"`
	AssetCategoryID          uint             `json:"asset_category_id"`
	AssetCategory            AssetCategoryDTO `json:"asset_category"`
}
