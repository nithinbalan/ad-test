package dtos

import (
	commonDTOs "jethur/grc/common/dtos"

	"time"
)

type RiskKriListDTO struct {
	ID                     uint    `json:"id"`
	KriReferenceNumber     string  `json:"kri_reference_number"`
	KriTitle               string  `json:"kri_title"`
	KRICategoryTitleEn     string  `json:"kri_category_title_en"`
	KRICategoryTitleAr     string  `json:"kri_category_title_ar"`
	OwnerUnitTitleEn       string  `json:"owner_unit_title_en" gorm:"column:owner_unit_title_en"`
	OwnerUnitTitleAr       string  `json:"owner_unit_title_ar" gorm:"column:owner_unit_title_ar"`
	OwnerFirstNameEn       string  `json:"owner_first_name_en"`
	OwnerFirstNameAr       string  `json:"owner_first_name_ar"`
	OwnerLastNameEn        string  `json:"owner_last_name_en"`
	OwnerLastNameAr        *string `json:"owner_last_name_ar"`
	ReviewFrequencyTitleEn string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr string  `json:"review_frequency_title_ar"`
	RiskLevelTitleEn       string  `json:"risk_level_title_en"`
	RiskLevelTitleAr       string  `json:"risk_level_title_ar"`
	RiskLevelColorCode     string  `json:"risk_level_color_code"`
	Description            *string `json:"description"`
	CalculationType        string  `json:"calculation_type"`
	LowThreshold           string  `json:"low_threshold"`
	MediumThreshold        string  `json:"medium_threshold"`
	HighThreshold          string  `json:"high_threshold"`
	CurrentValue           string  `json:"current_value"`
	LastReviewedDate       string  `json:"last_reviewed_date"`
	NextReviewDate         string  `json:"next_review_date"`
	LastReviewNote         string  `json:"last_review_note"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string  `json:"created_user_last_name_ar"`
}

type RiskKriListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type RiskKriListResponseDTO struct {
	Items      []*RiskKriListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskKriFormDTO struct {
	KriIDs []*uint `json:"kri_ids"`
}

type RiskKriParsedDates struct {
	TargetDate *time.Time
}

type RiskKriDetailDTO struct {
	ID          int                   `json:"id"`
	KriRegister KRIRegisterDetailsDTO `json:"kri_register"`
}
