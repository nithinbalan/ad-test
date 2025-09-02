package dtos

import commonDTOs "jethur/grc/common/dtos"

type ArticleLibraryDTO struct {
	ID              uint   `json:"id"`
	TitleEn         string `json:"title_en"`
	TitleAr         string `json:"title_ar"`
	ImagePath       string `json:"image_path"`
	CategoryID      uint   `json:"category_id"`
	CategoryTitleEn string `json:"category_title_en"`
	CategoryTitleAr string `json:"category_title_ar"`
	CreatedBy       string `json:"created_by"`
	CreatedAt       string `json:"created_at"`
}

type ArticleLibraryListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type ArticleLibraryListResponseDTO struct {
	Items      []*ArticleLibraryDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ArticleCategory struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
