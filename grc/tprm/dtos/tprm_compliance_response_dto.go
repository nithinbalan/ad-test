package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type TPRMComplianceResponseListDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Score     int    `json:"score"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TPRMComplianceResponseFormDTO struct {
	TitleEn   string `json:"title_en" binding:"required"`
	TitleAr   string `json:"title_ar" binding:"required"`
	Score     int    `json:"score"`
	ColorCode string `json:"color_code" binding:"required"`
}

type TPRMComplianceResponseDetailDTO struct {
	ID        uint       `json:"id"`
	TitleEn   string     `json:"title_en"`
	TitleAr   string     `json:"title_ar"`
	Score     int        `json:"score"`
	ColorCode string     `json:"color_code"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type TPRMComplianceResponseListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type TPRMComplianceResponseListResponseDTO struct {
	Items      []*TPRMComplianceResponseListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}
