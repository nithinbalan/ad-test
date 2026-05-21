package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type KRIRegisterListDTO struct {
	ID                              uint     `json:"id"`
	ReferenceNumber                 *string  `json:"reference_number"`
	KRICategoryTitleEn              string   `json:"kri_category_title_en"`
	KRICategoryTitleAr              string   `json:"kri_category_title_ar"`
	ResponsibleUnitTitleEn          string   `json:"responsible_unit_title_en" gorm:"column:responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string   `json:"responsible_unit_title_ar" gorm:"column:responsible_unit_title_ar"`
	ResponsibleUserID               uint     `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string   `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string   `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string   `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string   `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string  `json:"responsible_user_initial"`
	ResponsibleUserColor            *string  `json:"responsible_user_color"`
	ResponsibleUserEmail            *string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string  `json:"responsible_user_profile_image_path"`
	ReviewFrequencyTitleEn          string   `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          string   `json:"review_frequency_title_ar"`
	RiskLevelTitleEn                *string  `json:"risk_level_title_en"`
	RiskLevelTitleAr                *string  `json:"risk_level_title_ar"`
	RiskLevelColorCode              *string  `json:"risk_level_color_code"`
	UnitTitleEn                     *string  `json:"unit_title_en"`
	UnitTitleAr                     *string  `json:"unit_title_ar"`
	UnitCodeEn                      *string  `json:"unit_code_en"`
	UnitCodeAr                      *string  `json:"unit_code_ar"`
	UnitColorCode                   *string  `json:"unit_color_code"`
	LevelOfMeasurementID            *uint    `json:"level_of_measurement_id"`
	LevelOfMeasurementTitleEn       *string  `json:"level_of_measurement_title_en"`
	LevelOfMeasurementTitleAr       *string  `json:"level_of_measurement_title_ar"`
	RiskResponsibleUnitID           *uint    `json:"risk_responsible_unit_id"`
	RiskResponsibleUnitTitleEn      *string  `json:"risk_responsible_unit_title_en"`
	RiskResponsibleUnitTitleAr      *string  `json:"risk_responsible_unit_title_ar"`
	KRISubCategoryID                *uint    `json:"kri_sub_category_id"`
	KRISubCategoryTitleEn           *string  `json:"kri_sub_category_title_en"`
	KRISubCategoryTitleAr           *string  `json:"kri_sub_category_title_ar"`
	AnnualTarget                    *float64 `json:"annual_target"`
	ReviewTarget                    *float64 `json:"review_target"`
	DistributionStatusID            *uint    `json:"distribution_status_id"`
	DistributionStatusTitleEn       *string  `json:"distribution_status_title_en"`
	DistributionStatusTitleAr       *string  `json:"distribution_status_title_ar"`
	Title                           string   `json:"title"`
	Description                     *string  `json:"description"`
	CalculationType                 string   `json:"calculation_type"`
	LowThreshold                    *int     `json:"low_threshold"`
	MediumThreshold                 *int     `json:"medium_threshold"`
	HighThreshold                   *int     `json:"high_threshold"`
	CurrentValue                    *float64 `json:"current_value"`
	LastReviewedDate                *string  `json:"last_reviewed_date"`
	NextReviewDate                  string   `json:"next_review_date"`
	LastReviewNote                  *string  `json:"last_review_note"`
	MinimumSusceptibilityPercent    *float64 `json:"minimum_susceptibility_percent"`
	MaximumSusceptibilityPercent    *float64 `json:"maximum_susceptibility_percent"`
	MaximumTolerablePercent         *float64 `json:"maximum_tolerable_percent"`
	IsLowerBetter                   *bool    `json:"is_lower_better"`
	IsHigherBetter                  *bool    `json:"is_higher_better"`
	CreatedAt                       string   `json:"created_at"`
	UpdatedAt                       string   `json:"updated_at"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           *string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              *string  `json:"created_user_initial"`
	CreatedUserColor                *string  `json:"created_user_color"`
	CreatedUserEmail                *string  `json:"created_user_email"`
	CreatedUserProfileImagePath     *string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
	WorkflowStatusID                *uint    `json:"workflow_status_id"`
	WorkflowStatusTitleEn           *string  `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr           *string  `json:"workflow_status_title_ar"`
	WorkflowStatusType              *string  `json:"workflow_status_type"`
	WorkflowStatusColorCode         *string  `json:"workflow_status_color_code"`
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
	Export               bool
	DateFrom             string
	DateTo               string
}

type KRIRegisterListResponseDTO struct {
	Items      []*KRIRegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KRIRegisterFormDTO struct {
	KRICategoryID                uint                     `json:"kri_category_id" binding:"required"`
	ResponsibleUserID            uint                     `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID            uint                     `json:"review_frequency_id" binding:"required" gorm:"default:not null"`
	UnitID                       uint                     `json:"unit_id" binding:"required"`
	LevelOfMeasurementID         *uint                    `json:"level_of_measurement_id"`
	RiskResponsibleUnitID        *uint                    `json:"risk_responsible_unit_id"`
	KRISubCategoryID             *uint                    `json:"kri_sub_category_id"`
	AnnualTarget                 *float64                 `json:"annual_target"`
	ReviewTarget                 *float64                 `json:"review_target"`
	DistributionStatusID         *uint                    `json:"distribution_status_id"`
	Title                        string                   `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description                  string                   `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType              string                   `json:"calculation_type" binding:"required,oneof=direct ratio" gorm:"type:varchar(10);not null"`
	LowThreshold                 *int                     `json:"low_threshold" gorm:"type:int;null"`
	MediumThreshold              *int                     `json:"medium_threshold" gorm:"type:int;null"`
	CurrentValue                 *float64                 `json:"current_value" gorm:"type:decimal(18,2);null"`
	LastReviewedDate             string                   `json:"last_reviewed_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate               string                   `json:"next_review_date" binding:"required" gorm:"type:date;not null"`
	LastReviewNote               *string                  `json:"last_review_note" gorm:"type:nvarchar(max)"`
	MinimumSusceptibilityPercent *float64                 `json:"minimum_susceptibility_percent"`
	MaximumSusceptibilityPercent *float64                 `json:"maximum_susceptibility_percent"`
	MaximumTolerablePercent      *float64                 `json:"maximum_tolerable_percent"`
	IsLowerBetter                *bool                    `json:"is_lower_better"`
	IsHigherBetter               *bool                    `json:"is_higher_better"`
	FrequencyValues              []*FrequencyValueFormDTO `json:"frequency_values"`
}

type FrequencyValueFormDTO struct {
	Type        string  `json:"type" binding:"required"`
	TargetValue float64 `json:"target_value" binding:"required"`
}

type KRIFrequencyValueDTO struct {
	ID                  uint                       `json:"id"`
	KRIRegisterID       uint                       `json:"kri_register_id"`
	Year                int                        `json:"year"`
	MonthEn             string                     `json:"month_en"`
	MonthAr             string                     `json:"month_ar"`
	Type                string                     `json:"type"`
	TargetValue         float64                    `json:"target_value"`
	ActualValue         *float64                   `json:"actual_value"`
	AchievementValue    *float64                   `json:"achievement_value"`
	RiskToleranceStatus *KRIRiskToleranceStatusDTO `json:"risk_tolerance_status"`
	IsAchieved          bool                       `json:"is_achieved"`
	IsMissed            bool                       `json:"is_missed"`
}

type LevelOfMeasurementDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

type DistributionStatusDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
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
	LevelOfMeasurementID   *uint                             `json:"level_of_measurement_id"`
	LevelOfMeasurement     *LevelOfMeasurementDTO            `json:"level_of_measurement"`
	RiskResponsibleUnitID  *uint                             `json:"risk_responsible_unit_id"`
	RiskResponsibleUnit    *settingsDTOs.OrganizationUnitDTO `json:"risk_responsible_unit"`
	KRISubCategoryID       *uint                             `json:"kri_sub_category_id"`
	KRISubCategory         *KRISubCategoryDTO                `json:"kri_sub_category"`
	AnnualTarget           *float64                          `json:"annual_target"`
	ReviewTarget           *float64                          `json:"review_target"`
	DistributionStatusID   *uint                             `json:"distribution_status_id"`
	DistributionStatus     *DistributionStatusDTO            `json:"distribution_status"`
	Title                  string                            `json:"title"`
	Description            *string                           `json:"description"`
	CalculationType        string                            `json:"calculation_type"`
	LowThreshold           *int                              `json:"low_threshold"`
	MediumThreshold        *int                              `json:"medium_threshold"`
	HighThreshold          *int                              `json:"high_threshold"`
	CurrentValue           *float64                          `json:"current_value"`
	LastReviewedDateString string                            `json:"last_reviewed_date"`
	NextReviewDateString   string                            `json:"next_review_date"`
	LastReviewNote         *string                           `json:"last_review_note"`
	SubmittedByID          *uint                             `json:"submitted_by"`
	SubmittedAtString      *string                           `json:"submitted_at"`
	NextReviewLevel        *int                              `json:"next_review_level"`
	WorkflowStatusID       *uint                             `json:"workflow_status_id"`
	WorkflowStatus         *masterDTOs.WorkflowStatusDTO     `json:"workflow_status"`
	ApprovedByID           *uint                             `json:"approved_by"`
	ApprovedAtString       *string                           `json:"approved_at"`
	SubmittedBy            *organizationDTOs.UserMinimalDTO  `json:"submitted_by_user"`
	ApprovedBy             *organizationDTOs.UserMinimalDTO  `json:"approved_by_user"`
	// Stored percentage inputs
	MinimumSusceptibilityPercent *float64 `json:"minimum_susceptibility_percent"`
	MaximumSusceptibilityPercent *float64 `json:"maximum_susceptibility_percent"`
	MaximumTolerablePercent      *float64 `json:"maximum_tolerable_percent"`
	IsLowerBetter                *bool    `json:"is_lower_better"`
	IsHigherBetter               *bool    `json:"is_higher_better"`
	// Computed threshold values (percent * annual_target / 100)
	MinimumSusceptibility *float64                        `json:"minimum_susceptibility"`
	MaximumSusceptibility *float64                        `json:"maximum_susceptibility"`
	MaximumTolerable      *float64                        `json:"maximum_tolerable"`
	CreatedAt             time.Time                       `json:"created_at"`
	UpdatedAt             *time.Time                      `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser         organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
	ValueUpdates          []*KRIValueUpdateDetailsDTO     `json:"value_updates"`
	KRIFrequencyValues    []*KRIFrequencyValueDTO         `json:"frequency_values"`
}

type KRIRegisterParsedDates struct {
	LastReviewedDate *time.Time
	NextReviewDate   *time.Time
}

type KRITrendAnalysisDTO struct {
	ReviewedAtString string  `json:"reviewed_at"`
	Value            float64 `json:"value"`
}

// New DTO for frequency-based trend analysis
type KRIFrequencyTrendAnalysisDTO struct {
	ID               uint     `json:"id"`
	KRIRegisterID    uint     `json:"kri_register_id"`
	Year             int      `json:"year"`
	MonthEn          string   `json:"month_en"`
	MonthAr          string   `json:"month_ar"`
	Type             string   `json:"type"`
	TargetValue      float64  `json:"target_value"`
	ActualValue      *float64 `json:"actual_value"`
	AchievementValue *float64 `json:"achievement_value"`
	IsAchieved       bool     `json:"is_achieved"`
	IsMissed         bool     `json:"is_missed"`
	CreatedAt        string   `json:"created_at"`
	// Risk Tolerance Status details
	RiskToleranceStatusID  *uint   `json:"risk_tolerance_status_id"`
	RiskToleranceStatusEn  *string `json:"risk_tolerance_status_en"`
	RiskToleranceStatusAr  *string `json:"risk_tolerance_status_ar"`
	RiskToleranceColorCode *string `json:"risk_tolerance_color_code"`
	RiskToleranceType      *string `json:"risk_tolerance_type"`
}

type KRIFrequencyValueListDTO struct {
	ID               uint     `json:"id"`
	KRIRegisterID    uint     `json:"kri_register_id"`
	MonthEn          string   `json:"month_en"`
	MonthAr          string   `json:"month_ar"`
	Type             string   `json:"type"`
	TargetValue      float64  `json:"target_value"`
	ActualValue      *float64 `json:"actual_value"`
	AchievementValue *float64 `json:"achievement_value"`
	CreatedAt        string   `json:"created_at"`
}
