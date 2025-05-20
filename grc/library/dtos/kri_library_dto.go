package dtos

import commonDTOs "jethur/grc/common/dtos"

type KRILibraryDTO struct {
	ID          uint    `json:"id"`
	KRICategory string  `json:"kri_category"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	KRIUnit     *string `json:"kri_unit"`
}

type KRILibraryFormDTO struct {
	KRICategory string  `json:"kri_category" binding:"required,max=255"`
	Title       string  `json:"title" binding:"required,max=500"`
	Description *string `json:"description"`
	KRIUnit     *string `json:"kri_unit"`
}

type KRILibraryListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Category []string
	Page     int
	PageSize int
}

type KRILibraryListResponseDTO struct {
	Items      []*KRILibraryDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KRICategoryCountDTO struct {
	KRICategory string `json:"kri_category"`
	Count       int    `json:"count"`
}
