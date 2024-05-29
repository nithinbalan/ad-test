package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type KRIRegisterListDTO struct {
	ID                     uint    `json:"id"`
	ReferenceNumber        string  `json:"reference_number"`
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
	Title                  string  `json:"title"`
	Description            *string `json:"description"`
	CalculationType        string  `json:"calculation_type"`
	LowThreshold           string  `json:"low_threshold"`
	MediumThreshold        string  `json:"medium_threshold"`
	HighThreshold          string  `json:"high_threshold"`
	CurrentValue           string  `json:"current_value"`
	LastReviewedDate       string  `json:"last_reviewed_date"`
	NextReviewDate         string  `json:"next_review_date"`
	LastReviewNote         string  `json:"last_review_note"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  *string `json:"updated_user_last_name_ar"`
}

type KRIRegisterListQueryParamsDTO struct {
	Sort                 string
	SortBy               string
	Keyword              string
	Status               string
	Page                 int
	PageSize             int
	Filters              map[string][]string
	QuickFilter          []string
	LastReviewedDateFrom string
	LastReviewedDateTo   string
	NextReviewDateFrom   string
	NextReviewDateTo     string
}

type KRIRegisterListResponseDTO struct {
	Items      []*KRIRegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KRIRegisterFormDTO struct {
	KRICategoryID     uint    `json:"kri_category_id" binding:"required"`
	OwnerUnitID       uint    `json:"owner_unit_id" binding:"required"`
	OwnerID           uint    `json:"owner_id" binding:"required"`
	ReviewFrequencyID uint    `json:"review_frequency_id" binding:"required" gorm:"default:not null"`
	RiskLevelID       uint    `json:"risk_level_id" binding:"required" gorm:"default:not null"`
	Title             string  `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description       string  `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType   string  `json:"calculation_type" binding:"required,oneof=automatic manual" gorm:"type:varchar(10);not null"`
	LowThreshold      *string `json:"low_threshold" gorm:"type:varchar(455);null"`
	MediumThreshold   *string `json:"medium_threshold" gorm:"type:varchar(455);null"`
	HighThreshold     *string `json:"high_threshold" gorm:"type:varchar(455);null"`
	CurrentValue      *string `json:"current_value" gorm:"type:varchar(455);null"`
	LastReviewedDate  string  `json:"last_reviewed_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate    string  `json:"next_review_date" binding:"required" gorm:"type:date;not null"`
	LastReviewNote    *string `json:"last_review_note" gorm:"type:nvarchar(max)"`
}

type KRIRegisterDetailsDTO struct {
	ID               uint                              `json:"id"`
	ReferenceNumber  string                            `json:"reference_number"`
	KRICategory      KRICategoryDTO                    `json:"kri_category"`
	OwnerUnit        *settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	Owner            *organizationDTOs.UserMinimalDTO  `json:"owner"`
	ReviewFrequency  masterDTOs.ReviewFrequencyDTO     `json:"review_frequency"`
	RiskLevel        RiskLevelDTO                      `json:"risk_level"`
	Title            string                            `json:"title"`
	Description      *string                           `json:"description"`
	CalculationType  string                            `json:"calculation_type"`
	LowThreshold     *string                           `json:"low_threshold"`
	MediumThreshold  *string                           `json:"medium_threshold"`
	HighThreshold    *string                           `json:"high_threshold"`
	CurrentValue     *string                           `json:"current_value"`
	LastReviewedDate time.Time                         `json:"last_reviewed_date"`
	NextReviewDate   time.Time                         `json:"next_review_date"`
	LastReviewNote   *string                           `json:"last_review_note"`
	CreatedAt        time.Time                         `json:"created_at"`
	UpdatedAt        *time.Time                        `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO   `json:"created_by"`
	UpdatedByUser    organizationDTOs.UserMinimalDTO   `json:"updated_by"`
}

type KRIRegisterParsedDates struct {
	LastReviewedDate *time.Time
	NextReviewDate   *time.Time
}
