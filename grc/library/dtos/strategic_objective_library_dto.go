package dtos

import commonDTOs "jethur/grc/common/dtos"

type StrategicObjectiveLibraryDTO struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type StrategicObjectiveLibraryListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}

type StrategicObjectiveLibraryListResponseDTO struct {
	Items      []*StrategicObjectiveLibraryDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type StrategicCategoryCountDTO struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}
