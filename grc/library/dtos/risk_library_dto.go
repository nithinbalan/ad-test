package dtos

import commonDTOs "jethur/grc/common/dtos"

type RiskLibraryDTO struct {
	ID          uint    `json:"id"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Likelihood  *string `json:"likelihood"`
	Impact      *string `json:"impact"`
}

type RiskLibraryFormDTO struct {
	Category    string  `json:"category" binding:"required,max=255"`
	Title       string  `json:"title" binding:"required,max=500"`
	Description *string `json:"description"`
	Likelihood  *string `json:"likelihood"`
	Impact      *string `json:"impact"`
}

type RiskLibraryListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Category string
	Page     int
	PageSize int
}

type RiskLibraryListResponseDTO struct {
	Items      []*RiskLibraryDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskCategoryCountDTO struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}
