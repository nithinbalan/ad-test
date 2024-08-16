package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type KRIRegisterListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	KRICategoryTitleEn              string  `json:"kri_category_title_en"`
	KRICategoryTitleAr              string  `json:"kri_category_title_ar"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en" gorm:"column:responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar" gorm:"column:responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	ReviewFrequencyTitleEn          string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          string  `json:"review_frequency_title_ar"`
	RiskLevelTitleEn                string  `json:"risk_level_title_en"`
	RiskLevelTitleAr                string  `json:"risk_level_title_ar"`
	RiskLevelColorCode              string  `json:"risk_level_color_code"`
	UnitTitleEn                     string  `json:"unit_title_en"`
	UnitTitleAr                     string  `json:"unit_title_ar"`
	UnitCodeEn                      string  `json:"unit_code_en"`
	UnitCodeAr                      string  `json:"unit_code_ar"`
	UnitColorCode                   string  `json:"unit_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	CalculationType                 string  `json:"calculation_type"`
	LowThreshold                    *int    `json:"low_threshold"`
	MediumThreshold                 *int    `json:"medium_threshold"`
	HighThreshold                   *int    `json:"high_threshold"`
	CurrentValue                    int     `json:"current_value"`
	LastReviewedDate                string  `json:"last_reviewed_date"`
	NextReviewDate                  string  `json:"next_review_date"`
	LastReviewNote                  string  `json:"last_review_note"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       string  `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           *string `json:"created_user_last_name_ar"`
	CreatedUserInitial              *string `json:"created_user_initial"`
	CreatedUserColor                *string `json:"created_user_color"`
	CreatedUserEmail                *string `json:"created_user_email"`
	CreatedUserProfileImagePath     *string `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
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
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID uint    `json:"review_frequency_id" binding:"required" gorm:"default:not null"`
	RiskLevelID       *uint   `json:"risk_level_id" gorm:"default:null"`
	UnitID            uint    `json:"unit_id" binding:"required"`
	Title             string  `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description       string  `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType   string  `json:"calculation_type" binding:"required,oneof=automatic manual" gorm:"type:varchar(10);not null"`
	LowThreshold      *int    `json:"low_threshold" gorm:"type:int;null"`
	MediumThreshold   *int    `json:"medium_threshold" gorm:"type:int;null"`
	HighThreshold     *int    `json:"high_threshold" gorm:"type:int;null"`
	CurrentValue      *int    `json:"current_value" gorm:"type:int;null"`
	LastReviewedDate  string  `json:"last_reviewed_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate    string  `json:"next_review_date" binding:"required" gorm:"type:date;not null"`
	LastReviewNote    *string `json:"last_review_note" gorm:"type:nvarchar(max)"`
}

type KRIRegisterDetailsDTO struct {
	ID                     uint                              `json:"id"`
	ReferenceNumber        string                            `json:"reference_number"`
	KRICategory            KRICategoryDTO                    `json:"kri_category"`
	ResponsibleUnit        *settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser        *organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
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
	UpdatedByUser          organizationDTOs.UserMinimalDTO   `json:"updated_by_user"`
	ValueUpdates           []*KRIValueUpdateDetailsDTO       `json:"value_updates"`
}

type KRIRegisterParsedDates struct {
	LastReviewedDate *time.Time
	NextReviewDate   *time.Time
}
