package dtos

import commonDTOs "jethur/grc/common/dtos"

type KPILibraryDTO struct {
	ID          uint    `json:"id"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Unit        string  `json:"unit"`
	Frequency   string  `json:"frequency"`
}

type KPILibraryListQueryParamsDTO struct {
	Sort       string
	SortBy     string
	Keyword    string
	Page       int
	PageSize   int
	Export     bool
	Categories []string
}

type KPILibraryListResponseDTO struct {
	Items      []*KPILibraryDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPICategoryCountDTO struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}
