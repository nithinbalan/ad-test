package dtos

import "time"

type IntegrationDTO struct {
	ID                         int        `json:"id"`
	IntegrationCategoryID      uint       `json:"integration_category_id"`
	IntegrationCategoryTitleEn string     `json:"integration_category_title_en"`
	IntegrationCategoryTitleAr string     `json:"integration_category_title_ar"`
	TitleEn                    string     `json:"title_en"`
	TitleAr                    string     `json:"title_ar"`
	DescriptionEn              *string    `json:"description_en"`
	DescriptionAr              *string    `json:"description_ar"`
	LogoUrl                    *string    `json:"logo_url"`
	IsEnable                   bool       `json:"is_enable"`
	CreatedByID                uint       `json:"created_by_id"`
	UpdatedByID                *uint      `json:"updated_by_id"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  *time.Time `json:"updated_at"`
}

type IntegrationListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}
