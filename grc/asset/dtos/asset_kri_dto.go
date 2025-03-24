package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AssetKriListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	Title                           string  `json:"title"`
	KRIRegisterID                   uint    `json:"kri_register_id"`
	KRICategoryTitleEn              string  `json:"kri_category_title_en"`
	KRICategoryTitleAr              string  `json:"kri_category_title_ar"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_initial"`
	ResponsibleUserColor            *string `json:"responsible_color"`
	ResponsibleUserEmail            *string `json:"responsible_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_profile_image_path"`
	ReviewFrequencyTitleEn          string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          string  `json:"review_frequency_title_ar"`
	RiskLevelTitleEn                string  `json:"risk_level_title_en"`
	RiskLevelTitleAr                string  `json:"risk_level_title_ar"`
	RiskLevelColorCode              string  `json:"risk_level_color_code"`
	Description                     *string `json:"description"`
	CalculationType                 string  `json:"calculation_type"`
	LowThreshold                    int     `json:"low_threshold"`
	MediumThreshold                 int     `json:"medium_threshold"`
	HighThreshold                   int     `json:"high_threshold"`
	CurrentValue                    int     `json:"current_value"`
	LastReviewedDate                string  `json:"last_reviewed_date"`
	NextReviewDate                  string  `json:"next_review_date"`
	LastReviewNote                  string  `json:"last_review_note"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              *string `json:"created_user_initial"`
	CreatedUserColor                *string `json:"created_user_color"`
	CreatedUserEmail                *string `json:"created_user_email"`
	CreatedUserProfileImagePath     *string `json:"created_user_profile_image_path"`
}

type AssetKriListQueryParamsDTO struct {
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

type AssetKriListResponseDTO struct {
	Items      []*AssetKriListDTO           `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetKriFormDTO struct {
	KriIDs []*uint `json:"kri_ids"`
}

type AssetKriDetailDTO struct {
	ID            int                             `json:"id"`
	KriRegister   KRIRegisterDetailsDTO           `json:"kri_register"`
	CreatedByUser organizationDTOs.UserMinimalDTO `json:"created_by_user"`
}

type KRIRegisterDetailsDTO struct {
	ID                     uint                              `json:"id"`
	ReferenceNumber        string                            `json:"reference_number"`
	KRICategory            KRICategoryDTO                    `json:"kri_category"`
	ResponsibleUnit        *settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser        *organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ResponsibleUserID      uint                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequency        masterDTOs.ReviewFrequencyDTO     `json:"review_frequency"`
	RiskLevel              *RiskLevelDTO                     `json:"risk_level"`
	Unit                   *masterDTOs.UnitDTO               `json:"unit"`
	Title                  string                            `json:"title"`
	Description            *string                           `json:"description"`
	CalculationType        string                            `json:"calculation_type"`
	LowThreshold           *int                              `json:"low_threshold"`
	MediumThreshold        *int                              `json:"medium_threshold"`
	HighThreshold          *int                              `json:"high_threshold"`
	CurrentValue           *int                              `json:"current_value"`
	LastReviewedDateString string                            `json:"last_reviewed_date"`
	NextReviewDateString   string                            `json:"next_review_date"`
	LastReviewNote         *string                           `json:"last_review_note"`
	CreatedAt              time.Time                         `json:"created_at"`
	UpdatedAt              *time.Time                        `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO   `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO  `json:"updated_by_user"`
}

type KRICategoryDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type RiskLevelDTO struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
}
