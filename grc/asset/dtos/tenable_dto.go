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
	ReferenceNumber string    `json:"reference_number"`
	TitleEn         string    `json:"title_en"`
	Description     string    `json:"description"`
	StatusName      string    `json:"status_name"`
	Manufacturer    string    `json:"manufacturer"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TenableAssetListResponseDTO struct {
	Items      []TenableAssetDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
