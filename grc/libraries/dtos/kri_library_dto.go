package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type KRILibraryDTO struct {
	ID                     uint    `json:"id"`
	KRICategoryTitleEn     string  `json:"kri_category_title_en"`
	KRICategoryTitleAr     string  `json:"kri_category_title_ar"`
	ReviewFrequencyTitleEn string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr string  `json:"review_frequency_title_ar"`
	Title                  string  `json:"title"`
	Description            *string `json:"description"`
	CalculationType        string  `json:"calculation_type"`
	Source                 *string `json:"source"`
	Status                 string  `json:"status"`
}
type KRILibraryFormDTO struct {
	ID                uint    `json:"id"`
	KRICategoryID     uint    `json:"kri_category_id"`
	ReviewFrequencyID *uint   `json:"review_frequency_id"`
	Title             string  `json:"title"`
	Description       *string `json:"description"`
	CalculationType   string  `json:"calculation_type"`
	Source            *string `json:"source"`
	Status            string  `json:"status"`
}

type KRILibraryListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Status   string
	Page     int
	PageSize int
}

type KRILibraryListResponseDTO struct {
	Items      []*KRILibraryDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
