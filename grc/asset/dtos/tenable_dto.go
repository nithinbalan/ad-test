package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type TenableAssetListQueryParamsDTO struct {
	Page     int
	PageSize int
}

type TenableAssetDTO struct {
	ReferenceNumber    string    `json:"reference_number"`
	TitleEn            string    `json:"title_en"`
	Description        string    `json:"description"`
	StatusName         string    `json:"status_name"`
	Manufacturer       string    `json:"manufacturer"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	AssetCategoryID    uint      `json:"asset_category_id"`
	AssetSubCategoryID uint      `json:"asset_sub_category_id"`
	AssetStatusID      uint      `json:"asset_status_id"`
	ResponsibleUserID  uint      `json:"responsible_user_id"`
	AcquisitionDate    *string   `json:"acquisition_date" binding:"omitempty,datetime=2006-01-02"`
}

type TenableAssetListResponseDTO struct {
	Items      []TenableAssetDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
