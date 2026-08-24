package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type KPIRegisterListDTO struct {
	ID                              int      `json:"id"`
	ReferenceNumber                 *string  `json:"reference_number"`
	Title                           string   `json:"title"`
	Type                            string   `json:"type"`
	IsKRI                           bool     `json:"is_kri"`
	Description                     *string  `json:"description"`
	KPICategoryID                   uint     `json:"kpi_category_id"`
	KPICategoryTitleEn              string   `json:"kpi_category_title_en"`
	KPICategoryTitleAr              string   `json:"kpi_category_title_ar"`
	ResponsibleUserID               *uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string  `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string  `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string  `json:"responsible_user_email"`
	ResponsibleUserColor            *string  `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string  `json:"responsible_user_profile_image_path"`
	StatusID                        uint     `json:"status_id"`
	KPIStatusID                     uint     `json:"kpi_status_id"`
	KPIReviewStatusID               uint     `json:"kpi_review_status_id"`
	KPIStatusTitleEn                string   `json:"kpi_status_title_en"`
	KPIStatusTitleAr                string   `json:"kpi_status_title_ar"`
	KPIStatusColorCode              string   `json:"kpi_status_color_code"`
	KPIStatusType                   string   `json:"kpi_status_type"`
	StatusTitleEn                   string   `json:"status_title_en"`
	StatusTitleAr                   string   `json:"status_title_ar"`
	StatusColorCode                 string   `json:"status_color_code"`
	StatusType                      string   `json:"status_type"`
	KPIReviewStatusTitleEn          string   `json:"kpi_review_status_title_en"`
	KPIReviewStatusTitleAr          string   `json:"kpi_review_status_title_ar"`
	KPIReviewStatusColorCode        string   `json:"kpi_review_status_color_code"`
	KPIReviewStatusType             string   `json:"kpi_review_status_type"`
	ReviewFrequencyTitleEn          *string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string  `json:"review_frequency_title_ar"`
	ReviewFrequencyType             *string  `json:"review_frequency_type"`
	CalculationTypeTitleEn          *string  `json:"calculation_type_title_en"`
	CalculationTypeTitleAr          *string  `json:"calculation_type_title_ar"`
	CalculationTypeType             *string  `json:"calculation_type_type"`
	TargetValue                     *float64 `json:"target_value"`
	TargetDescription               *string  `json:"target_description"`
	PreviousValue                   *float64 `json:"previous_value"`
	Trend                           *float64 `json:"trend"`
	UnitTitleEn                     *string  `json:"unit_title_en"`
	UnitTitleAr                     *string  `json:"unit_title_ar"`
	UnitCodeEn                      *string  `json:"unit_code_en"`
	UnitCodeAr                      *string  `json:"unit_code_ar"`
	UnitColorCode                   *string  `json:"unit_color_code"`
	StartDate                       *string  `json:"start_date"`
	EndDate                         *string  `json:"end_date"`
	TargetDate                      *string  `json:"target_date"`
	Remarks                         *string  `json:"remarks"`
	Achievement                     *float64 `json:"achievement"`
	ActualValue                     *float64 `json:"actual_value"`
	KPIAchievementLogicID           uint     `json:"kpi_achievement_logic_id"`
	KPIAchievementLogicTitleEn      string   `json:"kpi_achievement_logic_title_en"`
	KPIAchievementLogicTitleAr      string   `json:"kpi_achievement_logic_title_ar"`
	KPIAchievementLogicType         string   `json:"kpi_achievement_logic_type"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string   `json:"created_user_last_name_ar"`
	CreatedUserInitial              string   `json:"created_user_initial"`
	CreatedUserColor                string   `json:"created_user_color"`
	CreatedUserEmail                string   `json:"created_user_email"`
	CreatedUserProfileImagePath     string   `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
	ReviewUserID                    *uint    `json:"review_user_id"`
	ReviewUserFirstNameEn           *string  `json:"review_user_first_name_en"`
	ReviewUserLastNameEn            *string  `json:"review_user_last_name_en"`
	ReviewUserFirstNameAr           *string  `json:"review_user_first_name_ar"`
	ReviewUserLastNameAr            *string  `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string  `json:"review_user_initial"`
	ReviewUserEmail                 *string  `json:"review_user_email"`
	ReviewUserColor                 *string  `json:"review_user_color"`
	ReviewUserProfileImagePath      *string  `json:"review_user_profile_image_path"`
}

type KPIRegisterFormDTO struct {
	Title             string              `json:"title" binding:"required"`
	IsKRI             bool                `json:"is_kri"`
	Description       *string             `json:"description"`
	KPICategoryID     uint                `json:"kpi_category_id" binding:"required"`
	ResponsibleUserID *uint               `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID *uint               `json:"review_frequency_id" binding:"required"`
	CalculationTypeID *uint               `json:"calculation_type_id" binding:"required"`
	DataInputs        []*InputDataFormDTO `json:"data_inputs"`
	Formula           *string             `json:"formula"`
	UnitID            *uint               `json:"unit_id" binding:"required"`

	TargetValue *float64   `json:"target_value" binding:"required"`
	TargetDate  *time.Time `json:"target_date"`

	StartDate             *string `json:"start_date"`
	EndDate               *string `json:"end_date"`
	Remarks               *string `json:"remarks" gorm:"type:nvarchar(max)"`
	KPIAchievementLogicID *uint   `json:"kpi_achievement_logic_id"`
	ReviewUserID          *uint   `json:"review_user_id"`
}

type InputDataFormDTO struct {
	Variable string `json:"variable"`
	Input    string `json:"input"`
}

type KPIRegisterListQueryParamsDTO struct {
	Sort          string
	SortBy        string
	Keyword       string
	Page          int
	PageSize      int
	Filters       map[string][]string
	QuickFilter   []string
	Export        bool
	CurrentUserID int
}

type KPIRegisterListResponseDTO struct {
	Items      []*KPIRegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIRegisterDetailDTO struct {
	ID                   uint                             `json:"id"`
	ReferenceNumber      *string                          `json:"reference_number"`
	Title                string                           `json:"title"`
	Type                 string                           `json:"type"`
	IsKRI                bool                             `json:"is_kri"`
	Description          *string                          `json:"description"`
	KPICategoryID        uint                             `json:"kpi_category_id"`
	KPICategory          *KPICategoryDTO                  `json:"kpi_category"`
	ResponsibleUserID    *uint                            `json:"responsible_user_id"`
	ResponsibleUser      *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	KPIStatusID          uint                             `json:"kpi_status_id"`
	KPIReviewStatusID    *uint                            `json:"kpi_review_status_id"`
	StatusID             *uint                            `json:"status_id"`
	KPIStatus            KPIStatusDTO                     `json:"kpi_status"`
	KPIReviewStatus      *KPIReviewStatusDTO              `json:"kpi_review_status"`
	Status               *masterDTOs.LifecycleStatusDTO   `json:"status"`
	ReviewFrequencyID    *uint                            `json:"review_frequency_id"`
	ReviewFrequency      *masterDTOs.ReviewFrequencyDTO   `json:"review_frequency"`
	CalculationTypeID    *uint                            `json:"calculation_type_id"`
	CalculationType      *masterDTOs.CalculationTypeDTO   `json:"calculation_type"`
	Formula              *string                          `json:"formula"`
	UnitID               *uint                            `json:"unit_id"`
	Unit                 *masterDTOs.UnitDTO              `json:"unit"`
	ActualValue          *float64                         `json:"actual_value"`
	TargetValue          *float64                         `json:"target_value"`
	TargetDescription    *string                          `json:"target_description"`
	Achievement          *float64                         `json:"achievement"`
	StartDateString      *string                          `json:"start_date"`
	EndDateString        *string                          `json:"end_date"`
	TargetDateString     *string                          `json:"target_date"`
	NextReviewDateString *string                          `json:"next_review_date"`
	LastReviewDateString *string                          `json:"last_review_date"`
	Remarks              *string                          `json:"remarks"`
	DataInputs           []*KPIRegisterDataInputDTO       `json:"data_inputs"`
	KPIAchievementLogic  *KPIAchievementLogicDTO          `json:"kpi_achievement_logic"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	ReviewUser           *organizationDTOs.UserMinimalDTO `json:"review_user"`
}

type KPIRegisterParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type KPIRegisterDataInputDTO struct {
	ID       uint   `json:"id"`
	Variable string `json:"variable"`
	Input    string `json:"input"`
}

type KPITrendAnalysisDTO struct {
	ReviewedAtString string  `json:"reviewed_at"`
	Value            float64 `json:"value"`
}

type KPIReviewStatusSummaryDTO struct {
	KPIReviewStatusDTO
	Count int64 `json:"count" gorm:"column:count"`
}

type KPICategorySummaryDTO struct {
	KPICategoryDTO
	Count          int64                        `json:"count"`
	Percentage     float64                      `json:"percentage"`
	ReviewStatuses []*KPIReviewStatusSummaryDTO `json:"review_statuses"`
}

type KPICategoryReviewCountDTO struct {
	CategoryID     uint  `json:"category_id" gorm:"column:category_id"`
	ReviewStatusID uint  `json:"review_status_id" gorm:"column:review_status_id"`
	Count          int64 `json:"count" gorm:"column:count"`
}

type KPIStatusSummaryDTO struct {
	KPIStatusDTO
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}
type KPIStatusSummaryResponseDTO struct {
	Statuses []*KPIStatusSummaryDTO `json:"statuses"`
}

type KPIOverallSummaryDTO struct {
	TotalKPIs          int64   `json:"total_kpis"`
	OverallAchievement float64 `json:"overall_achievement"`
}

type KPIAchievementLogicDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}
